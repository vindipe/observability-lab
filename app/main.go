package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	serviceName    = "orders-api"
	serviceVersion = "0.1.0"
)

var (
	logger      *slog.Logger
	redisClient *redis.Client
	tracer      trace.Tracer

	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "demo",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests handled by the demo service.",
		},
		[]string{"method", "route", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "demo",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1, 1.5, 2.5, 5},
		},
		[]string{"method", "route", "status"},
	)

	inFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "demo",
			Subsystem: "http",
			Name:      "in_flight_requests",
			Help:      "Current number of in-flight HTTP requests.",
		},
	)

	dependencyDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "demo",
			Subsystem: "dependency",
			Name:      "duration_seconds",
			Help:      "External dependency operation duration in seconds.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 0.75, 1, 2, 5},
		},
		[]string{"dependency", "operation", "status"},
	)

	dependencyErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "demo",
			Subsystem: "dependency",
			Name:      "errors_total",
			Help:      "Total number of dependency errors seen by the demo service.",
		},
		[]string{"dependency", "operation"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration, inFlightRequests, dependencyDuration, dependencyErrors)
}

func main() {
	ctx := context.Background()
	logger = initLogger()

	shutdownTracer, err := initTracer(ctx)
	if err != nil {
		logger.Error("failed to initialize OpenTelemetry tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			logger.Error("failed to shutdown tracer provider", "error", err)
		}
	}()

	redisClient = initRedisClient()
	tracer = otel.Tracer(serviceName)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", instrumentRoute("/", rootHandler))
	mux.HandleFunc("GET /health", instrumentRoute("/health", healthHandler))
	mux.HandleFunc("GET /api/orders", instrumentRoute("/api/orders", ordersHandler))
	mux.HandleFunc("GET /api/slow", instrumentRoute("/api/slow", slowHandler))
	mux.HandleFunc("GET /api/error", instrumentRoute("/api/error", errorHandler))
	mux.Handle("GET /metrics", promhttp.Handler())

	wrapped := otelhttp.NewHandler(mux, "http.server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))

	addr := ":" + getenv("APP_PORT", "8080")
	server := &http.Server{
		Addr:              addr,
		Handler:           wrapped,
		ReadHeaderTimeout: 3 * time.Second,
	}

	logger.Info("starting orders-api", "addr", addr, "service", serviceName, "version", serviceVersion)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("http server failed", "error", err)
		os.Exit(1)
	}
}

func initLogger() *slog.Logger {
	writers := []io.Writer{os.Stdout}
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		if err := os.MkdirAll(strings.TrimSuffix(logFile, "/app.log"), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "cannot create log directory: %v\n", err)
		} else if file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
			writers = append(writers, file)
		} else {
			fmt.Fprintf(os.Stderr, "cannot open log file %s: %v\n", logFile, err)
		}
	}

	handler := slog.NewJSONHandler(io.MultiWriter(writers...), &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(handler).With(
		"service", serviceName,
		"service_version", serviceVersion,
		"deployment_environment", getenv("DEPLOYMENT_ENVIRONMENT", "local"),
	)
}

func initTracer(ctx context.Context) (func(context.Context) error, error) {
	endpoint := getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "otel-collector:4317")
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	sampleRatio := parseFloat(getenv("OTEL_SAMPLE_RATIO", "1.0"), 1.0)
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			attribute.String("deployment.environment", getenv("DEPLOYMENT_ENVIRONMENT", "local")),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRatio))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp.Shutdown, nil
}

func initRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		logger.Warn("REDIS_ADDR not set: dependency calls will be simulated without Redis")
		return nil
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("redis not reachable at startup; app continues and records dependency errors", "addr", addr, "error", err)
	} else {
		logger.Info("redis dependency reachable", "addr", addr)
	}
	return client
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func instrumentRoute(route string, next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inFlightRequests.Inc()
		defer inFlightRequests.Dec()

		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(recorder, r)

		status := strconv.Itoa(recorder.status)
		duration := time.Since(start)
		httpRequests.WithLabelValues(r.Method, route, status).Inc()
		httpDuration.WithLabelValues(r.Method, route, status).Observe(duration.Seconds())

		requestLogger := logger.With(
			"method", r.Method,
			"route", route,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
		if sc := trace.SpanFromContext(r.Context()).SpanContext(); sc.IsValid() {
			requestLogger = requestLogger.With("trace_id", sc.TraceID().String(), "span_id", sc.SpanID().String())
		}
		if recorder.status >= 500 {
			requestLogger.Error("request completed")
		} else {
			requestLogger.Info("request completed")
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service":   serviceName,
		"endpoints": []string{"/health", "/api/orders", "/api/slow", "/api/error", "/metrics"},
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "health-check")
	defer span.End()

	redisStatus := "disabled"
	if redisClient != nil {
		pingCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
		defer cancel()
		if err := redisClient.Ping(pingCtx).Err(); err != nil {
			redisStatus = "degraded"
			span.RecordError(err)
			span.SetAttributes(attribute.String("redis.status", redisStatus))
		} else {
			redisStatus = "ok"
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": serviceName,
		"redis":   redisStatus,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "orders.list")
	defer span.End()

	sleepRandom(ctx, 50*time.Millisecond, 250*time.Millisecond)
	_ = recordRedisOperation(ctx, "incr", func(redisCtx context.Context) error {
		if redisClient == nil {
			return nil
		}
		return redisClient.Incr(redisCtx, "orders-api:orders_requests_total").Err()
	})

	orders := []map[string]any{
		{"id": "ord-1001", "status": "paid", "amount": 19.90},
		{"id": "ord-1002", "status": "preparing", "amount": 42.00},
		{"id": "ord-1003", "status": "shipped", "amount": 12.50},
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": orders})
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "external.dependency.slow-query")
	defer span.End()

	sleepRandom(ctx, 700*time.Millisecond, 1600*time.Millisecond)
	err := recordRedisOperation(ctx, "getset", func(redisCtx context.Context) error {
		if redisClient == nil {
			return nil
		}
		key := fmt.Sprintf("orders-api:slow:%d", rand.Intn(5))
		if err := redisClient.Set(redisCtx, key, time.Now().UTC().Format(time.RFC3339Nano), 30*time.Second).Err(); err != nil {
			return err
		}
		return redisClient.Get(redisCtx, key).Err()
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "dependency unavailable", "detail": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"note":   "simulated slow dependency call completed",
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "intentional.error")
	defer span.End()

	errorRate := parseFloat(getenv("DEMO_ERROR_RATE", "0.40"), 0.40)
	sleepRandom(ctx, 25*time.Millisecond, 100*time.Millisecond)
	if rand.Float64() < errorRate {
		err := errors.New("intentional demo error")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn("returning intentional error", "error_rate", errorRate)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "lucky request: no error this time"})
}

func recordRedisOperation(ctx context.Context, operation string, fn func(context.Context) error) error {
	ctx, span := tracer.Start(ctx, "redis."+operation, trace.WithAttributes(
		attribute.String("db.system", "redis"),
		attribute.String("db.operation", operation),
	))
	defer span.End()

	start := time.Now()
	redisCtx, cancel := context.WithTimeout(ctx, 900*time.Millisecond)
	defer cancel()

	err := fn(redisCtx)
	status := "ok"
	if err != nil {
		status = "error"
		dependencyErrors.WithLabelValues("redis", operation).Inc()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	dependencyDuration.WithLabelValues("redis", operation, status).Observe(time.Since(start).Seconds())
	return err
}

func sleepRandom(ctx context.Context, min, max time.Duration) {
	if max <= min {
		select {
		case <-time.After(min):
		case <-ctx.Done():
		}
		return
	}
	delta := rand.Int63n(int64(max - min))
	select {
	case <-time.After(min + time.Duration(delta)):
	case <-ctx.Done():
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("failed to write JSON response", "error", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseFloat(value string, fallback float64) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

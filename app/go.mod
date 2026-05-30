module github.com/vindipe/observability-lab/app

go 1.23

require (
	github.com/prometheus/client_golang v1.21.1
	github.com/redis/go-redis/v9 v9.7.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.59.0
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
)

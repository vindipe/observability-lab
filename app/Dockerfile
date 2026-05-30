FROM golang:1.23-bookworm AS build
WORKDIR /src

# Copy module metadata and source before resolving modules.
# In a freshly generated lab repository go.sum can be empty; `go mod tidy`
# materializes the required checksums from the actual imports before build.
COPY go.mod ./
COPY main.go ./
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/orders-api .

FROM debian:bookworm-slim
RUN useradd -u 10001 -r -s /usr/sbin/nologin appuser \
    && mkdir -p /var/log/observability-lab \
    && chown -R appuser:appuser /var/log/observability-lab
COPY --from=build /out/orders-api /usr/local/bin/orders-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/orders-api"]

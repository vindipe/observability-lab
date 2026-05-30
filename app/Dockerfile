FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/orders-api .

FROM debian:bookworm-slim
RUN useradd -u 10001 -r -s /usr/sbin/nologin appuser \
    && mkdir -p /var/log/observability-lab \
    && chown -R appuser:appuser /var/log/observability-lab
COPY --from=build /out/orders-api /usr/local/bin/orders-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/orders-api"]

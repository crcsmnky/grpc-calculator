FROM golang:1.18-buster as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GO111MODULE=on go build -v -o grpc-calculator github.com/crcsmnky/grpc-calculator/server

FROM debian:buster-slim

RUN set -x && \
    apt-get update && \
    apt-get install -y ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

ADD https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.60.0/otelcol-contrib_0.60.0_linux_amd64.tar.gz /otelcol-contrib_0.60.0.tar.gz
RUN tar zxf otelcol-contrib_0.60.0.tar.gz

ADD https://github.com/krallin/tini/releases/download/v0.19.0/tini /tini
RUN chmod +x /tini

COPY entrypoint.sh config.yaml ./

COPY --from=build /app/grpc-calculator /app/grpc-calculator

# COPY --from=datadog/serverless-init:beta2 /datadog-init /app/datadog-init

# ENV DD_SERVICE=grpc-calculator-otel
# ENV DD_ENV=verily-test
# ENV DD_VERSION=1

# ENV DD_OTLP_CONFIG_RECEIVER_PROTOCOLS_GRPC_ENDPOINT=localhost:4317
ENV OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317

# ENTRYPOINT ["/app/datadog-init"]
# CMD ["/app/grpc-calculator"]

ENTRYPOINT ["/tini", "--", "/entrypoint.sh"]

FROM golang:1.18-buster as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -v -o grpc-calculator github.com/crcsmnky/grpc-calculator/server

FROM debian:buster-slim

RUN set -x && \
    apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /app/grpc-calculator /app/grpc-calculator
COPY --from=datadog/serverless-init:beta3 /datadog-init /app/datadog-init

ENV DD_SERVICE=grpc-calculator-ddtrace
ENV DD_ENV=verily-test
ENV DD_VERSION=1

ENTRYPOINT ["/app/datadog-init"]
CMD ["/app/grpc-calculator"]
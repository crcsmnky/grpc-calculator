#!/usr/bin/env bash

set -e

/app/grpc-calculator &
/otelcol-contrib --config config.yaml &

wait -n

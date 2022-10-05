#!/usr/bin/env bash

for((count=0; count<50; count++)); do
  for op in "ADD" "SUBTRACT" "MULTIPLY" "DIVIDE"; do
    rando=$((RANDOM%100))
    # result=$(grpcurl -proto proto/calculator.proto \
    #   -d '{"first_operand":"'$rando'", "second_operand":"2.0", "operation":"'$op'"}' \
    #   -H "Authorization: Bearer $(gcloud --project parikhs-experiments auth print-identity-token)" \
    #   grpc-calculator-otel-y7xk6uygoa-uc.a.run.app:443 Calculator.Calculate | jq '.result')

    result=$(grpcurl -proto proto/calculator.proto -plaintext \
      -d '{"first_operand":"'$rando'", "second_operand":"2.0", "operation":"'$op'"}' \
      35.235.87.165:3000 Calculator.Calculate | jq '.result')

    echo "$count: $rando $op 2.0 = $result"
    sleep 2
  done
done
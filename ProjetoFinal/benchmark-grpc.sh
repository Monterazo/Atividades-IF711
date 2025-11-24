#!/bin/bash
# Script de Benchmark para gRPC
# Uso: ./benchmark-grpc.sh [expressao] [clientes] [requisicoes]

EXPRESSION=${1:-"((4+3)*2)/5"}
CLIENTS=${2:-10}
REQUESTS=${3:-100}

echo ""
echo "============================================"
echo "  BENCHMARK gRPC - Calculadora Distribuída"
echo "============================================"
echo "  Expressão: $EXPRESSION"
echo "  Clientes:  $CLIENTS"
echo "  Requests:  $REQUESTS"
echo "============================================"
echo ""

echo "Compilando benchmark..."
go build -o bin/grpc_benchmark cmd/benchmark/grpc_benchmark.go

if [ $? -ne 0 ]; then
    echo "Erro ao compilar benchmark!"
    exit 1
fi

echo "Executando benchmark..."
echo ""

./bin/grpc_benchmark -expr="$EXPRESSION" -clients=$CLIENTS -reqs=$REQUESTS

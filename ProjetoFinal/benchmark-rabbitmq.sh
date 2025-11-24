#!/bin/bash
# Script de Benchmark para RabbitMQ
# Uso: ./benchmark-rabbitmq.sh [expressao] [clientes] [requisicoes]

EXPRESSION=${1:-"((4+3)*2)/5"}
CLIENTS=${2:-10}
REQUESTS=${3:-100}

echo ""
echo "================================================"
echo "  BENCHMARK RabbitMQ - Calculadora Distribuída"
echo "================================================"
echo "  Expressão: $EXPRESSION"
echo "  Clientes:  $CLIENTS"
echo "  Requests:  $REQUESTS"
echo "================================================"
echo ""

echo "Compilando benchmark..."
go build -o bin/rabbitmq_benchmark cmd/benchmark/rabbitmq_benchmark.go

if [ $? -ne 0 ]; then
    echo "Erro ao compilar benchmark!"
    exit 1
fi

echo "Executando benchmark..."
echo ""

./bin/rabbitmq_benchmark -expr="$EXPRESSION" -clients=$CLIENTS -reqs=$REQUESTS

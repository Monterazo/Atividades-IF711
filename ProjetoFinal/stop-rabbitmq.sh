#!/bin/bash

echo "Parando sistema RabbitMQ..."

# Lê PIDs salvos
if [ -f /tmp/rabbitmq_pids.txt ]; then
    PIDS=$(cat /tmp/rabbitmq_pids.txt)
    for PID in $PIDS; do
        echo "  Parando PID $PID..."
        kill $PID 2>/dev/null || kill -9 $PID 2>/dev/null
    done
    rm /tmp/rabbitmq_pids.txt
    echo "✅ Processos salvos parados!"
else
    echo "⚠️  Arquivo de PIDs não encontrado, tentando por nome..."
fi

# Fallback: mata por nome
echo "  Buscando processos restantes..."
pkill -f rabbitmq_add_server 2>/dev/null
pkill -f rabbitmq_sub_server 2>/dev/null
pkill -f rabbitmq_mult_server 2>/dev/null
pkill -f rabbitmq_div_server 2>/dev/null
pkill -f rabbitmq_dispatcher 2>/dev/null
pkill -f rabbitmq_client 2>/dev/null

echo "✅ Sistema RabbitMQ parado!"

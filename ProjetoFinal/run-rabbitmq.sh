#!/bin/bash

echo "Iniciando sistema RabbitMQ..."
echo "IMPORTANTE: Certifique-se de que o RabbitMQ está rodando!"
echo ""

# Compila se necessário
if [ ! -f "bin/rabbitmq_client.exe" ]; then
    echo "Compilando binários RabbitMQ..."
    make build-rabbitmq
fi

# Inicia servidores em background
echo "Iniciando servidores de operação..."
./bin/rabbitmq_add_server.exe > /tmp/rabbitmq_add.log 2>&1 &
ADD_PID=$!
echo "  [ADD] PID $ADD_PID"

./bin/rabbitmq_sub_server.exe > /tmp/rabbitmq_sub.log 2>&1 &
SUB_PID=$!
echo "  [SUB] PID $SUB_PID"

./bin/rabbitmq_mult_server.exe > /tmp/rabbitmq_mult.log 2>&1 &
MULT_PID=$!
echo "  [MULT] PID $MULT_PID"

./bin/rabbitmq_div_server.exe > /tmp/rabbitmq_div.log 2>&1 &
DIV_PID=$!
echo "  [DIV] PID $DIV_PID"

# Aguarda servidores iniciarem
echo ""
echo "Aguardando servidores iniciarem..."
sleep 2

# Inicia dispatcher
echo "Iniciando dispatcher..."
./bin/rabbitmq_dispatcher.exe > /tmp/rabbitmq_dispatcher.log 2>&1 &
DISPATCHER_PID=$!
echo "  [DISPATCHER] PID $DISPATCHER_PID"

# Aguarda dispatcher iniciar
sleep 2

# Salva PIDs para cleanup
echo "$ADD_PID $SUB_PID $MULT_PID $DIV_PID $DISPATCHER_PID" > /tmp/rabbitmq_pids.txt

echo ""
echo "✅ Sistema iniciado!"
echo "📋 Logs disponíveis em /tmp/rabbitmq_*.log"
echo ""
echo "Para parar todos os processos, execute:"
echo "  ./stop-rabbitmq.sh"
echo ""
echo "Iniciando cliente..."
echo "================================"
echo ""

# Inicia cliente (foreground)
./bin/rabbitmq_client.exe

# Cleanup ao sair do cliente
echo ""
echo "Cliente encerrado. Parando servidores..."
if [ -f /tmp/rabbitmq_pids.txt ]; then
    PIDS=$(cat /tmp/rabbitmq_pids.txt)
    for PID in $PIDS; do
        kill $PID 2>/dev/null
    done
    rm /tmp/rabbitmq_pids.txt
fi

echo "✅ Sistema encerrado!"

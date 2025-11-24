#!/bin/bash
# Script para parar todos os servidores gRPC

echo "🛑 Parando servidores gRPC..."
echo ""

if [ ! -f ".grpc_pids" ]; then
    echo "⚠️  Arquivo de PIDs não encontrado."
    echo "   Tentando matar processos pelo nome..."
    pkill -f grpc_add_server
    pkill -f grpc_sub_server
    pkill -f grpc_mult_server
    pkill -f grpc_div_server
    pkill -f grpc_dispatcher
    echo "✅ Processos finalizados."
    exit 0
fi

# Lê PIDs do arquivo e mata os processos
while read pid; do
    if ps -p $pid > /dev/null 2>&1; then
        echo "   Matando processo PID: $pid"
        kill $pid 2>/dev/null
    fi
done < .grpc_pids

# Aguarda um pouco
sleep 1

# Força a parada se ainda estiverem rodando
while read pid; do
    if ps -p $pid > /dev/null 2>&1; then
        echo "   Forçando parada do PID: $pid"
        kill -9 $pid 2>/dev/null
    fi
done < .grpc_pids

# Remove arquivo de PIDs
rm -f .grpc_pids

echo ""
echo "✅ Todos os servidores gRPC foram parados!"
echo ""

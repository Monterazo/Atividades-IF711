#!/bin/bash
# Script para iniciar todos os servidores gRPC em background

echo "🚀 Iniciando servidores gRPC..."
echo ""

# Verifica se os binários existem
if [ ! -f "bin/grpc_add_server.exe" ]; then
    echo "❌ Binários não encontrados. Compilando..."
    make build-grpc
    if [ $? -ne 0 ]; then
        echo "❌ Erro ao compilar. Abortando."
        exit 1
    fi
fi

# Cria diretório de logs se não existir
mkdir -p logs

# Inicia servidores de operação
echo "📦 Iniciando servidor de adição..."
./bin/grpc_add_server.exe > logs/grpc_add_server.log 2>&1 &
ADD_PID=$!
echo "   PID: $ADD_PID"

echo "📦 Iniciando servidor de subtração..."
./bin/grpc_sub_server.exe > logs/grpc_sub_server.log 2>&1 &
SUB_PID=$!
echo "   PID: $SUB_PID"

echo "📦 Iniciando servidor de multiplicação..."
./bin/grpc_mult_server.exe > logs/grpc_mult_server.log 2>&1 &
MULT_PID=$!
echo "   PID: $MULT_PID"

echo "📦 Iniciando servidor de divisão..."
./bin/grpc_div_server.exe > logs/grpc_div_server.log 2>&1 &
DIV_PID=$!
echo "   PID: $DIV_PID"

echo ""
echo "⏳ Aguardando servidores iniciarem..."
sleep 2

# Inicia dispatcher
echo "🎯 Iniciando dispatcher..."
./bin/grpc_dispatcher.exe > logs/grpc_dispatcher.log 2>&1 &
DISP_PID=$!
echo "   PID: $DISP_PID"

echo ""
echo "⏳ Aguardando dispatcher iniciar..."
sleep 2

# Salva PIDs em arquivo
echo "$ADD_PID" > .grpc_pids
echo "$SUB_PID" >> .grpc_pids
echo "$MULT_PID" >> .grpc_pids
echo "$DIV_PID" >> .grpc_pids
echo "$DISP_PID" >> .grpc_pids

echo ""
echo "✅ Todos os servidores gRPC estão rodando!"
echo ""
echo "📋 Processos:"
echo "   Add Server:    PID $ADD_PID"
echo "   Sub Server:    PID $SUB_PID"
echo "   Mult Server:   PID $MULT_PID"
echo "   Div Server:    PID $DIV_PID"
echo "   Dispatcher:    PID $DISP_PID"
echo ""
echo "📁 Logs salvos em: logs/"
echo "🛑 Para parar os servidores: ./stop-grpc.sh"
echo ""
echo "🚀 Você pode agora executar o benchmark:"
echo "   make benchmark-grpc EXPR=\"2+2\" CLIENTS=5 REQS=50"
echo ""

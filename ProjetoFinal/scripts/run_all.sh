#!/bin/bash
# Script para executar o sistema completo de calculadora distribuída

echo "======================================"
echo "Iniciando Sistema de Calculadora gRPC"
echo "======================================"
echo ""

# Função para limpar processos ao sair
cleanup() {
    echo ""
    echo "======================================"
    echo "Encerrando servidores..."
    echo "======================================"
    kill $ADD_PID $SUB_PID $MULT_PID $DIV_PID $DISP_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

# Iniciar servidores de operação em background
echo "Iniciando servidores de operação..."
./bin/grpc_add_server &
ADD_PID=$!
echo "  ✓ Servidor de Adição (PID: $ADD_PID)"

./bin/grpc_sub_server &
SUB_PID=$!
echo "  ✓ Servidor de Subtração (PID: $SUB_PID)"

./bin/grpc_mult_server &
MULT_PID=$!
echo "  ✓ Servidor de Multiplicação (PID: $MULT_PID)"

./bin/grpc_div_server &
DIV_PID=$!
echo "  ✓ Servidor de Divisão (PID: $DIV_PID)"

# Aguardar servidores iniciarem
echo ""
echo "Aguardando servidores iniciarem..."
sleep 2

# Iniciar dispatcher em background
echo "Iniciando dispatcher..."
./bin/grpc_dispatcher &
DISP_PID=$!
echo "  ✓ Dispatcher (PID: $DISP_PID)"

# Aguardar dispatcher iniciar
sleep 2

# Iniciar cliente (foreground)
echo ""
echo "======================================"
echo "Iniciando cliente..."
echo "======================================"
echo ""
./bin/grpc_client

# Limpar ao finalizar
cleanup

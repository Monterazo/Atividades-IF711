# ⚡ Benchmark - Guia Rápido

## 🚀 Comandos Principais

### Benchmark gRPC
```bash
# 1. Iniciar servidores
./run-grpc.sh

# 2. Executar benchmark
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100

# 3. Parar servidores
./stop-grpc.sh
```

### Benchmark RabbitMQ
```bash
# 1. Iniciar servidores
./run-rabbitmq.sh

# 2. Executar benchmark
make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100

# 3. Parar servidores
./stop-rabbitmq.sh
```

## 📝 Parâmetros

```bash
make benchmark-grpc EXPR="<expressão>" CLIENTS=<número> REQS=<número>
```

- **EXPR**: Expressão matemática (ex: "2+2", "((4+3)*2)/5")
- **CLIENTS**: Número de clientes simultâneos (ex: 5, 10, 50)
- **REQS**: Requisições por cliente (ex: 50, 100, 200)

## 🧪 Exemplos Prontos

```bash
# Teste simples
make benchmark-grpc EXPR="2+2" CLIENTS=5 REQS=50

# Teste médio
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100

# Teste pesado
make benchmark-grpc EXPR="((10+5)*3-7)/2" CLIENTS=50 REQS=200

# Teste de erro
make benchmark-grpc EXPR="10/0" CLIENTS=5 REQS=10
```

## 📊 O que o Benchmark Mostra

```
📈 Requisições:
   Total:        1000
   Sucesso:      995 (99.50%)
   Falhas:       5 (0.50%)

⏱️  Latência:
   Mínima:       12.5ms
   Média:        45.3ms
   Máxima:       156.7ms
   P50:          42.1ms    ← 50% das requisições
   P95:          89.4ms    ← 95% das requisições
   P99:          124.2ms   ← 99% das requisições

🚀 Desempenho:
   Duração total:    10.5s
   Throughput:       94.76 req/s
```

## 🛠 Comandos Úteis

```bash
# Compilar benchmarks
make build-benchmark

# Executar diretamente (mais controle)
./bin/grpc_benchmark.exe -expr="2+2" -clients=5 -reqs=50 -v

# Ver logs
tail -f logs/grpc_dispatcher.log

# Listar processos rodando
ps aux | grep grpc
ps aux | grep rabbitmq

# Matar processos manualmente
./stop-grpc.sh
./stop-rabbitmq.sh
```

## 📖 Documentação Completa

Para instruções detalhadas, cenários avançados e troubleshooting:
- [BENCHMARK_GUIDE.md](BENCHMARK_GUIDE.md) - Guia completo
- [README.md](README.md) - Visão geral do projeto

## 💡 Dicas Rápidas

1. **Sempre inicie os servidores antes do benchmark**
2. **Use valores menores para testes rápidos** (CLIENTS=5, REQS=50)
3. **Modo verboso (-v) é útil para debug**
4. **Pare os servidores após os testes**
5. **Compare gRPC vs RabbitMQ com mesmos parâmetros**

## 🔥 Comparação Rápida

Execute os mesmos testes em ambos para comparar:

```bash
# Teste 1: gRPC
./run-grpc.sh
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100
./stop-grpc.sh

# Teste 2: RabbitMQ
./run-rabbitmq.sh
make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100
./stop-rabbitmq.sh
```

Compare:
- Latência média
- Throughput (req/s)
- P95 / P99
- Taxa de sucesso

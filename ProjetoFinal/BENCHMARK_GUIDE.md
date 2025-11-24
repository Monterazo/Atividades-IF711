# 📊 Guia de Benchmark - Calculadora Distribuída

Este guia explica como usar as ferramentas de benchmark automatizado para testar o desempenho do sistema com múltiplos clientes simultâneos.

## 🎯 Objetivo

As ferramentas de benchmark permitem:
- Testar o sistema com múltiplos clientes concorrentes
- Medir latência (tempo de resposta) de requisições
- Calcular throughput (requisições por segundo)
- Analisar percentis de latência (P50, P95, P99)
- Comparar desempenho entre gRPC e RabbitMQ

## 📦 Ferramentas Disponíveis

### 1. Benchmark gRPC
Testa o sistema usando comunicação gRPC síncrona.

### 2. Benchmark RabbitMQ
Testa o sistema usando comunicação assíncrona via RabbitMQ.

## 🚀 Como Usar

### Método 1: Scripts Prontos (Mais Fácil)

#### Windows:
```bash
# Benchmark gRPC
benchmark-grpc.bat "((4+3)*2)/5" 10 100

# Benchmark RabbitMQ
benchmark-rabbitmq.bat "((4+3)*2)/5" 10 100
```

#### Linux/Mac:
```bash
# Tornar scripts executáveis (apenas uma vez)
chmod +x benchmark-grpc.sh
chmod +x benchmark-rabbitmq.sh

# Benchmark gRPC
./benchmark-grpc.sh "((4+3)*2)/5" 10 100

# Benchmark RabbitMQ
./benchmark-rabbitmq.sh "((4+3)*2)/5" 10 100
```

**Parâmetros:**
1. Expressão matemática (padrão: `((4+3)*2)/5`)
2. Número de clientes simultâneos (padrão: 10)
3. Número de requisições por cliente (padrão: 100)

### Método 2: Makefile

```bash
# Benchmark gRPC com parâmetros padrão
make benchmark-grpc

# Benchmark gRPC customizado
make benchmark-grpc EXPR="10+20*3" CLIENTS=50 REQS=200

# Benchmark RabbitMQ com parâmetros padrão
make benchmark-rabbitmq

# Benchmark RabbitMQ customizado
make benchmark-rabbitmq EXPR="(15-5)/2" CLIENTS=20 REQS=500
```

### Método 3: Executáveis Diretos (Mais Controle)

Primeiro, compile os benchmarks:
```bash
make build-benchmark
```

#### gRPC:
```bash
bin/grpc_benchmark.exe -expr="((4+3)*2)/5" -clients=10 -reqs=100 -timeout=30000 -v
```

#### RabbitMQ:
```bash
bin/rabbitmq_benchmark.exe -expr="((4+3)*2)/5" -clients=10 -reqs=100 -timeout=30000 -v
```

**Flags disponíveis:**
- `-expr`: Expressão matemática a ser testada
- `-clients`: Número de clientes simultâneos
- `-reqs`: Número de requisições por cliente
- `-timeout`: Timeout em milissegundos (padrão: 30000)
- `-v`: Modo verboso (mostra cada requisição)
- `-dispatcher`: Endereço do dispatcher (apenas gRPC, padrão: localhost:50051)
- `-url`: URL do RabbitMQ (apenas RabbitMQ, padrão: amqp://guest:guest@localhost:5672/)

## 📋 Pré-requisitos

### Para Benchmark gRPC:
1. Iniciar todos os servidores gRPC:
   ```bash
   # Em terminais separados ou use o script run-grpc.sh
   bin/grpc_add_server.exe
   bin/grpc_sub_server.exe
   bin/grpc_mult_server.exe
   bin/grpc_div_server.exe
   bin/grpc_dispatcher.exe
   ```

### Para Benchmark RabbitMQ:
1. Garantir que o RabbitMQ está rodando:
   ```bash
   # Windows
   rabbitmq-server

   # Linux
   sudo systemctl start rabbitmq-server

   # Mac
   brew services start rabbitmq
   ```

2. Iniciar todos os servidores RabbitMQ:
   ```bash
   # Em terminais separados
   bin/rabbitmq_add_server.exe
   bin/rabbitmq_sub_server.exe
   bin/rabbitmq_mult_server.exe
   bin/rabbitmq_div_server.exe
   bin/rabbitmq_dispatcher.exe
   ```

## 📊 Interpretando os Resultados

O benchmark exibe estatísticas detalhadas:

```
============================================================
📊 RESULTADOS DO BENCHMARK
============================================================

📈 Requisições:
   Total:        1000
   Sucesso:      995 (99.50%)
   Falhas:       5 (0.50%)

⏱️  Latência:
   Mínima:       12.5ms
   Média:        45.3ms
   Máxima:       156.7ms
   P50:          42.1ms
   P95:          89.4ms
   P99:          124.2ms

🚀 Desempenho:
   Duração total:    10.5s
   Throughput:       94.76 req/s

============================================================
```

### Métricas Explicadas:

- **Total**: Total de requisições enviadas
- **Sucesso**: Requisições que obtiveram resposta correta
- **Falhas**: Requisições que falharam ou deram timeout
- **Latência Mínima**: Menor tempo de resposta observado
- **Latência Média**: Média de todos os tempos de resposta
- **Latência Máxima**: Maior tempo de resposta observado
- **P50 (Mediana)**: 50% das requisições foram mais rápidas que este valor
- **P95**: 95% das requisições foram mais rápidas que este valor
- **P99**: 99% das requisições foram mais rápidas que este valor
- **Duração total**: Tempo total do benchmark
- **Throughput**: Requisições bem-sucedidas por segundo

## 🧪 Cenários de Teste Recomendados

### 1. Teste de Carga Leve
```bash
make benchmark-grpc EXPR="2+2" CLIENTS=5 REQS=50
make benchmark-rabbitmq EXPR="2+2" CLIENTS=5 REQS=50
```

### 2. Teste de Carga Média
```bash
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=20 REQS=100
make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=20 REQS=100
```

### 3. Teste de Carga Pesada
```bash
make benchmark-grpc EXPR="((10+5)*3-7)/2" CLIENTS=50 REQS=200
make benchmark-rabbitmq EXPR="((10+5)*3-7)/2" CLIENTS=50 REQS=200
```

### 4. Teste de Stress
```bash
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=100 REQS=500
make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=100 REQS=500
```

### 5. Teste de Expressões Complexas
```bash
make benchmark-grpc EXPR="((15-5)*2+(10/2))*3" CLIENTS=30 REQS=150
make benchmark-rabbitmq EXPR="((15-5)*2+(10/2))*3" CLIENTS=30 REQS=150
```

### 6. Teste de Erros (Divisão por Zero)
```bash
make benchmark-grpc EXPR="10/0" CLIENTS=10 REQS=50
make benchmark-rabbitmq EXPR="10/0" CLIENTS=10 REQS=50
```

## 📈 Comparação gRPC vs RabbitMQ

Para comparar o desempenho, execute os mesmos testes em ambos os sistemas:

```bash
# Teste 1: gRPC
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=50 REQS=200 > results_grpc.txt

# Teste 2: RabbitMQ
make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=50 REQS=200 > results_rabbitmq.txt
```

Compare as métricas:
- **Latência**: Qual sistema responde mais rápido?
- **Throughput**: Qual processa mais requisições por segundo?
- **Taxa de Sucesso**: Qual tem menos falhas?
- **Comportamento sob carga**: Como cada um se comporta com muitos clientes?

## 🔍 Modo Verboso

Para debug e análise detalhada, use o modo verboso:

```bash
bin/grpc_benchmark.exe -expr="((4+3)*2)/5" -clients=2 -reqs=5 -v
```

Isso mostrará cada requisição individual:
```
✅ [Cliente 0] Conectado ao dispatcher
✅ [Cliente 0 | Req 0] Resultado: 2.800000 (tempo: 45ms)
✅ [Cliente 0 | Req 1] Resultado: 2.800000 (tempo: 42ms)
...
```

## 🐛 Troubleshooting

### Erro: "Falha ao conectar ao dispatcher"
- **gRPC**: Verifique se o dispatcher está rodando em `localhost:50051`
- **RabbitMQ**: Verifique se o RabbitMQ está rodando

### Erro: "Timeout aguardando resposta"
- Aumente o timeout: `-timeout=60000` (60 segundos)
- Verifique se todos os servidores de operação estão rodando
- Reduza o número de clientes/requisições simultâneas

### Taxa de Falhas Alta
- Verifique logs dos servidores para identificar erros
- Reduza a carga (menos clientes ou requisições)
- Verifique recursos do sistema (CPU, memória)

### Desempenho Baixo
- Verifique se há outros processos consumindo recursos
- Teste com expressões mais simples
- Monitore uso de CPU e memória durante o teste

## 💡 Dicas

1. **Warm-up**: Execute um teste pequeno antes do benchmark principal para "aquecer" o sistema
2. **Múltiplas execuções**: Execute cada teste 3-5 vezes e tire a média dos resultados
3. **Isolamento**: Feche outros programas durante os testes
4. **Monitoramento**: Use ferramentas de monitoramento do sistema (Task Manager, top, htop)
5. **Logs**: Analise os logs dos servidores para identificar gargalos

## 📝 Exemplo de Relatório

```markdown
# Resultados de Benchmark - Calculadora Distribuída

## Configuração do Teste
- Expressão: ((4+3)*2)/5
- Clientes: 50
- Requisições por cliente: 200
- Total de requisições: 10.000

## Resultados gRPC
- Latência média: 45.3ms
- Throughput: 94.76 req/s
- Taxa de sucesso: 99.5%
- P95: 89.4ms

## Resultados RabbitMQ
- Latência média: 78.6ms
- Throughput: 63.82 req/s
- Taxa de sucesso: 98.2%
- P95: 156.7ms

## Conclusão
O gRPC apresentou melhor desempenho em latência e throughput,
enquanto o RabbitMQ oferece maior resiliência e desacoplamento.
```

## 🎓 Próximos Passos

Após executar os benchmarks:
1. Documente os resultados
2. Compare as duas arquiteturas
3. Identifique gargalos
4. Analise trade-offs entre as abordagens
5. Prepare relatório de avaliação de desempenho

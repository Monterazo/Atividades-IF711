# 📘 **Arquitetura do Sistema Distribuído – MOM (RabbitMQ) e RPC (gRPC)**

## 📌 **Visão Geral**

Este projeto implementa dois estilos de comunicação distribuída, seguindo as instruções fornecidas no arquivo especificacao.txt:

**Arquitetura MOM (Message-Oriented Middleware)**
- Baseada em RabbitMQ ✅ **(IMPLEMENTADO)**
- Responsável por comunicação assíncrona, baseada em mensagens.

**Arquitetura RPC (Remote Procedure Call)**
- Baseada em gRPC ✅ **(IMPLEMENTADO)**
- Comunicação síncrona e tipada entre processos distribuídos.

O objetivo é implementar uma calculadora distribuída capaz de avaliar expressões matemáticas complexas enviadas pelo cliente. As expressões são quebradas em etapas pelo Dispatcher e enviadas aos servidores especializados (Add, Sub, Mult, Div).

O projeto exige ainda um relatório comparativo de desempenho entre as duas abordagens.

---

## 🚀 **Quick Start**

### Pré-requisitos
- Go 1.21+
- **Para gRPC:** Protocol Buffers Compiler (protoc)
- **Para RabbitMQ:** RabbitMQ Server
- Git

### Instalação e Execução Rápida - gRPC

```bash
# 1. Instalar protoc (Windows - PowerShell como Admin)
powershell -ExecutionPolicy Bypass -File scripts\install_protoc.ps1

# 2. Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Compilar e executar
make run-all
```

### Instalação e Execução Rápida - RabbitMQ

```bash
# 1. Instalar e iniciar RabbitMQ
# Windows: choco install rabbitmq && rabbitmq-server
# Linux: sudo apt install rabbitmq-server && sudo systemctl start rabbitmq-server
# macOS: brew install rabbitmq && brew services start rabbitmq

# 2. Compilar e executar
make run-all-rabbitmq
```

### Documentação Detalhada
- 📖 [SETUP.md](SETUP.md) - Configuração completa do ambiente
- 📋 [INSTRUCOES.md](INSTRUCOES.md) - Instruções detalhadas gRPC
- 🐰 [INSTRUCOES_RABBITMQ.md](INSTRUCOES_RABBITMQ.md) - Instruções detalhadas RabbitMQ
- 📡 [especificacao.txt](especificacao.txt) - Especificação do projeto

---

## 🧱 **1. Arquitetura Lógica Comum (Core Layer)**

Independente de MQTT ou gRPC, o sistema possui a mesma estrutura conceitual.

### ✔ **Componentes:**

**1. Cliente**
- Interface CLI simples.
- Envia expressão matemática.
- Aguarda resposta (com timeout).

**2. Dispatcher**

Peça central do sistema.

Responsável por:
- Parsing da expressão (Shunting Yard → RPN).
- Decomposição em steps (operações atômicas).
- Envio das operações aos servidores especializados.
- Reagrupamento das respostas.
- Tratamento de timeouts, erros e fluxo de execução.

**3. Servidores Especializados**

Cada servidor implementa apenas UMA operação:

| Servidor | Operação | Exemplo |
|----------|----------|---------|
| AddServer | Soma | 4 + 3 → 7 |
| SubServer | Subtração | 7 - 2 → 5 |
| MultServer | Multiplicação | 5 * 3 → 15 |
| DivServer | Divisão | 14 / 5 → 2.8 |

## 🧩 **2. Modelo de Dados Padronizado**

**Request de Expressão**
```json
{
  "expression_id": "expr_abc123",
  "expression": "((4+3)*2)/5",
  "deadline_ms": 30000
}
```

**Response de Expressão**
```json
{
  "expression_id": "expr_abc123",
  "result": 2.8,
  "error": null
}
```

**OperationRequest**
```json
{
  "expression_id": "expr_abc123",
  "step_id": "expr_abc123_step1",
  "operation": "add",
  "numbers": [4, 3],
  "deadline_ms": 5000
}
```

**OperationResponse**
```json
{
  "expression_id": "expr_abc123",
  "step_id": "expr_abc123_step1",
  "result": 7,
  "error": null
}
```

**ErrorInfo**
```json
{
  "code": "DIV_BY_ZERO",
  "message": "Cannot divide by zero."
}
```

## 🎯 **3. Parsing e Execução (RPN)**

**Exemplo:** `((4 + 3) * 2) / 5`

**Passos:**

| Step | Operação | Entrada | Saída |
|------|----------|---------|-------|
| 1 | add | 4, 3 | 7 |
| 2 | multiply | 7, 2 | 14 |
| 3 | divide | 14, 5 | 2.8 |

Dispatcher coordena exatamente estes passos.

## 📡 **4. Arquitetura MOM (RabbitMQ)**

### 📊 **4.1 Diagrama**
```
Cliente
   │
   ▼
[RabbitMQ Broker] ←→ Dispatcher ←→ Servidores (Add/Sub/Mult/Div)
```

### 📚 **4.2 Filas RabbitMQ (Padronizadas)**

**Requests:**
- `calculator.requests`

**Responses:**
- `calculator.responses`

**Operações:**
- `operations.add`
- `operations.subtract`
- `operations.multiply`
- `operations.divide`

**Resultados dos servidores:**
- `operations.results`

### 🔁 **4.3 Fluxo de Execução RabbitMQ**

1. Cliente → `calculator.requests`.
2. Dispatcher consome, faz parsing.
3. Para cada step:
   - Publica OperationRequest na fila correta (operations.add, etc.).
4. Servidor especializado:
   - Processa
   - Publica em `operations.results`.
5. Dispatcher coleta, ordena e monta o resultado final.
6. Publica resultado em `calculator.responses`.

### 🛠 **Melhorias aplicadas à arquitetura**

- ✔ Utilizado RabbitMQ como broker MOM
- ✔ Padronizado JSON como serialização oficial
- ✔ Estruturados IDs (expressionId, stepId)
- ✔ Separado core da implementação RabbitMQ
- ✔ Filas duráveis para garantir persistência de mensagens
- ✔ Documentação revisada e padronizada

## ⚡ **5. Arquitetura RPC (gRPC)**

Agora a versão distribuída via chamadas diretas RPC.

### 📊 **5.1 Diagrama gRPC**
```
Cliente → Dispatcher → Servidores Específicos
```

É um pipeline síncrono com context timeout.

### 📜 **5.2 Definição .proto (padrão oficial do projeto)**

**Serviço Cliente → Dispatcher**
```protobuf
service CalculatorService {
  rpc Calculate(ExpressionRequest) returns (ExpressionResponse);
}
```

**Serviço Dispatcher → Servidores**
```protobuf
service OperationService {
  rpc Execute(OperationRequest) returns (OperationResponse);
}
```

**Mensagens**
```protobuf
message ExpressionRequest {
  string expression_id = 1;
  string expression = 2;
  int64 deadline_ms = 3;
}

message ExpressionResponse {
  string expression_id = 1;
  double result = 2;
  ErrorInfo error = 3;
}

message OperationRequest {
  string expression_id = 1;
  string step_id = 2;
  string operation = 3;
  repeated double numbers = 4;
  int64 deadline_ms = 5;
}

message OperationResponse {
  string expression_id = 1;
  string step_id = 2;
  double result = 3;
  ErrorInfo error = 4;
}
```

## 🚀 **6. Fluxo RPC Completo**

1. **Cliente chama:**
   - `CalculatorService.Calculate()`

2. **Dispatcher:**
   - Faz parsing
   - Converte para RPN
   - Para cada step:
     - Chama o servidor certo via: `OperationService.Execute()`

3. **Se qualquer step falhar:**
   - Erro é devolvido imediatamente.

4. **Se tudo der certo:**
   - Dispatcher monta o resultado final e retorna ao cliente.

## 🏛 **7. Estrutura de Pastas Implementada**
```
/ProjetoFinal
│
├── cmd/
│   ├── rabbitmq_dispatcher/    ✅ IMPLEMENTADO
│   ├── rabbitmq_add_server/    ✅ IMPLEMENTADO
│   ├── rabbitmq_sub_server/    ✅ IMPLEMENTADO
│   ├── rabbitmq_mult_server/   ✅ IMPLEMENTADO
│   ├── rabbitmq_div_server/    ✅ IMPLEMENTADO
│   ├── rabbitmq_client/        ✅ IMPLEMENTADO
│   ├── grpc_dispatcher/        ✅ IMPLEMENTADO
│   ├── grpc_add_server/        ✅ IMPLEMENTADO
│   ├── grpc_sub_server/        ✅ IMPLEMENTADO
│   ├── grpc_mult_server/       ✅ IMPLEMENTADO
│   ├── grpc_div_server/        ✅ IMPLEMENTADO
│   └── grpc_client/            ✅ IMPLEMENTADO
│
├── internal/
│   ├── core/        # Parsing, modelos e regras comuns ✅
│   ├── rabbitmq/    # Implementação RabbitMQ ✅
│   └── grpc/        # Implementação gRPC ✅
│
└── proto/           # Definições Protocol Buffers ✅
```

## 📊 **8. Benchmark e Comparação de Desempenho**

### ⚡ Quick Start - Executar Benchmarks

**Benchmark gRPC:**
```bash
./run-grpc.sh                                           # Iniciar servidores
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100
./stop-grpc.sh                                          # Parar servidores
```

**Benchmark RabbitMQ:**
```bash
./run-rabbitmq.sh                                       # Iniciar servidores
make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100
./stop-rabbitmq.sh                                      # Parar servidores
```

**Documentação completa:** Ver [BENCHMARK_GUIDE.md](BENCHMARK_GUIDE.md)

### 📈 Métricas Coletadas

Os benchmarks automatizados medem:
- ✅ Latência (mínima, média, máxima, P50, P95, P99)
- ✅ Throughput (requisições por segundo)
- ✅ Taxa de sucesso/falha
- ✅ Duração total do teste
- ✅ Comportamento sob carga concorrente

### 🧪 Cenários de Teste Disponíveis

```bash
# Carga leve
make benchmark-grpc EXPR="2+2" CLIENTS=5 REQS=50

# Carga média
make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=20 REQS=100

# Carga pesada
make benchmark-grpc EXPR="((10+5)*3-7)/2" CLIENTS=50 REQS=200

# Teste de erro
make benchmark-grpc EXPR="10/0" CLIENTS=10 REQS=50
```

## 🎯 **9. Conclusão**

Este documento e implementação unificam:

- ✔ A especificação oficial do projeto
- ✔ A arquitetura MOM com RabbitMQ (IMPLEMENTADA)
- ✔ A arquitetura RPC com gRPC (IMPLEMENTADA)
- ✔ Código compartilhado no pacote `internal/core`
- ✔ Duas implementações completas e funcionais
- ✔ Documentação detalhada para ambas as abordagens
- ✔ Scripts de build e execução automatizados
- ✔ Pronto para testes de benchmark e comparação de desempenho
- ✔ Pronto para apresentação, entrega e avaliação

---

## 🚀 **Próximos Passos Recomendados**

1. **Executar ambas as implementações** para validar funcionamento
2. **Implementar benchmarks** para comparar desempenho
3. **Coletar métricas** de latência, throughput, uso de CPU/memória
4. **Criar relatório comparativo** entre MOM e RPC
5. **Testar cenários de falha** e recuperação
6. **Documentar observações** e lições aprendidas
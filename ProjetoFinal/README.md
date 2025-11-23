# 📘 **Arquitetura do Sistema Distribuído – MOM (MQTT) e RPC (gRPC)**

## 📌 **Visão Geral**

Este projeto implementa dois estilos de comunicação distribuída, seguindo as instruções fornecidas no arquivo especificacao.txt:

**Arquitetura MOM (Message-Oriented Middleware)**
- Baseada em MQTT
- Responsável por comunicação assíncrona, baseada em mensagens.

**Arquitetura RPC (Remote Procedure Call)**
- Baseada em gRPC ✅ **(IMPLEMENTADO)**
- Comunicação síncrona e tipada entre processos distribuídos.

O objetivo é implementar uma calculadora distribuída capaz de avaliar expressões matemáticas complexas enviadas pelo cliente. As expressões são quebradas em etapas pelo Dispatcher e enviadas aos servidores especializados (Add, Sub, Mult, Div).

O projeto exige ainda um relatório comparativo de desempenho entre as duas abordagens.

---

## 🚀 **Quick Start - gRPC**

### Pré-requisitos
- Go 1.21+
- Protocol Buffers Compiler (protoc)
- Git

### Instalação e Execução Rápida

```bash
# 1. Instalar protoc (Windows - PowerShell como Admin)
powershell -ExecutionPolicy Bypass -File scripts\install_protoc.ps1

# 2. Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Compilar o projeto
powershell -ExecutionPolicy Bypass -File scripts\build.ps1

# 4. Executar o sistema
powershell -ExecutionPolicy Bypass -File scripts\run.ps1
```

### Documentação Detalhada
- 📖 [SETUP.md](SETUP.md) - Configuração completa do ambiente
- 📋 [INSTRUCOES.md](INSTRUCOES.md) - Instruções detalhadas de execução
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

## 📡 **4. Arquitetura MOM (MQTT)**

### 📊 **4.1 Diagrama**
```
Cliente
   │
   ▼
[MQTT Broker] ←→ Dispatcher ←→ Servidores (Add/Sub/Mult/Div)
```

### 📚 **4.2 Tópicos MQTT (Padronizados)**

**Requests:**
- `calculator/requests`

**Responses:**
- `calculator/responses`

**Operações:**
- `operations/add`
- `operations/subtract`
- `operations/multiply`
- `operations/divide`

**Resultados dos servidores:**
- `operations/results`

### 🔁 **4.3 Fluxo de Execução MQTT**

1. Cliente → `calculator/requests`.
2. Dispatcher consome, faz parsing.
3. Para cada step:
   - Publica OperationRequest no tópico correto (operations/add, etc.).
4. Servidor especializado:
   - Processa
   - Publica em `operations/results`.
5. Dispatcher coleta, ordena e monta o resultado final.
6. Publica resultado em `calculator/responses`.

### 🛠 **Melhorias aplicadas à arquitetura original**

- ✔ Removido conceito de Connection Pool MQTT (não necessário)
- ✔ Padronizado JSON como serialização oficial
- ✔ Mantido MessagePack como opcional
- ✔ Estruturados IDs (expressionId, stepId)
- ✔ Separado core da implementação MQTT
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

## 🏛 **7. Estrutura de Pastas Recomendada**
```
/calculator-distributed
│
├── cmd/
│   ├── mqtt_dispatcher/
│   ├── mqtt_add_server/
│   ├── mqtt_sub_server/
│   ├── mqtt_mult_server/
│   ├── mqtt_div_server/
│   ├── mqtt_client/
│   ├── grpc_dispatcher/
│   ├── grpc_add_server/
│   ├── grpc_sub_server/
│   ├── grpc_mult_server/
│   ├── grpc_div_server/
│   └── grpc_client/
│
├── internal/
│   ├── core/        # Parsing, modelos e regras comuns
│   ├── mqtt/        # Implementação MQTT
│   └── grpc/        # Implementação gRPC
│
└── bench/
    ├── benchmark_mqtt.go
    ├── benchmark_grpc.go
    └── results.md
```

## 📊 **8. Benchmark e Comparação de Desempenho**

**Métricas a serem coletadas:**
- Latência total (p50, p95)
- Throughput (req/s)
- Uso de memória
- Uso de CPU
- Taxa de falhas
- Impacto de concorrência

**Cenários recomendados:**
- 10.000 expressões simples
- 5.000 expressões complexas
- 50 clientes simultâneos
- Casos de erro (divisão por zero)
- Testes com latência artificial

## 🎯 **9. Conclusão**

Este documento unifica:

- ✔ A especificação oficial
- ✔ A arquitetura MOM do colega (corrigida)
- ✔ A arquitetura RPC gRPC (padronizada)
- ✔ As boas práticas da disciplina
- ✔ Uma estrutura de repositório profissional
- ✔ Pronto para apresentação, entrega e avaliação
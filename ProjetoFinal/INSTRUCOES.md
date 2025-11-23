# 📋 **Instruções de Execução - Calculadora Distribuída gRPC**

## 🔧 **Pré-requisitos**

Antes de executar o projeto, certifique-se de ter instalado:

1. **Go 1.21+** - [Download](https://golang.org/dl/)
2. **Protocol Buffers Compiler (protoc)** - [Download](https://grpc.io/docs/protoc-installation/)
3. **Git** (para clonar o repositório)

### Instalação do protoc no Windows:
```bash
# Baixar protoc do GitHub
# https://github.com/protocolbuffers/protobuf/releases

# Adicionar ao PATH do sistema
# Exemplo: C:\protoc\bin
```

### Instalação dos plugins Go para protoc:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 📦 **Instalação**

### 1. Clonar o repositório (se ainda não fez)
```bash
cd C:\Users\mikae\Documents\GitHub\Atividades-IF711\ProjetoFinal
```

### 2. Baixar dependências
```bash
go mod download
go mod tidy
```

### 3. Gerar código a partir do .proto
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/calculator.proto
```

Ou use o Makefile:
```bash
make proto
```

### 4. Compilar os binários
```bash
make build
```

Isso irá criar os seguintes executáveis em `bin/`:
- `grpc_add_server.exe` - Servidor de adição
- `grpc_sub_server.exe` - Servidor de subtração
- `grpc_mult_server.exe` - Servidor de multiplicação
- `grpc_div_server.exe` - Servidor de divisão
- `grpc_dispatcher.exe` - Dispatcher central
- `grpc_client.exe` - Cliente CLI

## 🚀 **Execução**

### Opção 1: Execução Manual (Recomendado para debug)

Abra **5 terminais diferentes** e execute em ordem:

**Terminal 1 - Servidor Add:**
```bash
bin/grpc_add_server.exe
```

**Terminal 2 - Servidor Subtract:**
```bash
bin/grpc_sub_server.exe
```

**Terminal 3 - Servidor Multiply:**
```bash
bin/grpc_mult_server.exe
```

**Terminal 4 - Servidor Divide:**
```bash
bin/grpc_div_server.exe
```

**Terminal 5 - Dispatcher:**
```bash
bin/grpc_dispatcher.exe
```

**Terminal 6 - Cliente:**
```bash
bin/grpc_client.exe
```

### Opção 2: Execução Automatizada

Use o script de execução:
```bash
make run-all
```

Ou manualmente com scripts:

**Windows (PowerShell):**
```powershell
# Iniciar servidores
Start-Process -NoNewWindow .\bin\grpc_add_server.exe
Start-Process -NoNewWindow .\bin\grpc_sub_server.exe
Start-Process -NoNewWindow .\bin\grpc_mult_server.exe
Start-Process -NoNewWindow .\bin\grpc_div_server.exe

# Aguardar inicialização
Start-Sleep -Seconds 2

# Iniciar dispatcher
Start-Process -NoNewWindow .\bin\grpc_dispatcher.exe

# Aguardar inicialização
Start-Sleep -Seconds 2

# Iniciar cliente
.\bin\grpc_client.exe
```

## 💻 **Usando o Cliente**

Após iniciar o cliente, você verá:
```
Cliente Calculadora gRPC
========================
Conectado ao dispatcher em localhost:50051

Digite uma expressão matemática (ou 'sair' para encerrar):
Exemplos: ((4+3)*2)/5, 10+20*3, (15-5)/2

>
```

### Exemplos de Expressões:

```
> ((4+3)*2)/5
✅ Resultado: ((4+3)*2)/5 = 2.800000
⏱️  Tempo de execução: 15ms

> 10+20*3
✅ Resultado: 10+20*3 = 70.000000
⏱️  Tempo de execução: 12ms

> (15-5)/2
✅ Resultado: (15-5)/2 = 5.000000
⏱️  Tempo de execução: 10ms

> 10/0
❌ Erro: [DIV_BY_ZERO] divisão por zero
```

Para sair do cliente, digite `sair` ou `exit`.

## 🏗️ **Arquitetura**

```
Cliente (porta CLI)
    │
    ▼
Dispatcher (porta 50051)
    │
    ├──> AddServer (porta 50052)
    ├──> SubServer (porta 50053)
    ├──> MultServer (porta 50054)
    └──> DivServer (porta 50055)
```

### Portas Utilizadas:
- **50051** - Dispatcher (CalculatorService)
- **50052** - Add Server (OperationService)
- **50053** - Subtract Server (OperationService)
- **50054** - Multiply Server (OperationService)
- **50055** - Divide Server (OperationService)

## 🔍 **Fluxo de Execução**

1. **Cliente** envia expressão para o **Dispatcher**
2. **Dispatcher** faz parsing da expressão (Shunting Yard → RPN)
3. **Dispatcher** decompõe em steps e envia para servidores especializados
4. Cada **Servidor** executa sua operação e retorna resultado
5. **Dispatcher** agrupa resultados e retorna ao **Cliente**

## 🐛 **Troubleshooting**

### Erro: "Falha ao conectar ao dispatcher"
- Verifique se o dispatcher está rodando
- Confirme que a porta 50051 está livre

### Erro: "Falha ao conectar aos servidores"
- Verifique se todos os 4 servidores de operação estão rodando
- Confirme que as portas 50052-50055 estão livres

### Erro: "protoc: command not found"
- Instale o Protocol Buffers Compiler
- Adicione ao PATH do sistema

### Erro ao compilar .proto
```bash
# Reinstale os plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 🧪 **Testes**

### Teste de Operações Básicas:
```
> 5+3          # Esperado: 8
> 10-4         # Esperado: 6
> 6*7          # Esperado: 42
> 15/3         # Esperado: 5
```

### Teste de Expressões Complexas:
```
> ((4+3)*2)/5       # Esperado: 2.8
> 10+20*3           # Esperado: 70
> (15-5)/2          # Esperado: 5
> ((8+2)*3)/6       # Esperado: 5
```

### Teste de Erros:
```
> 10/0              # Esperado: Erro DIV_BY_ZERO
> abc+def           # Esperado: Erro PARSE_ERROR
```

## 📊 **Logs**

Os servidores geram logs detalhados:

**Dispatcher:**
```
Recebida expressão: ((4+3)*2)/5 (ID: expr_1)
Expressão parseada em 3 steps
Executando step 0: add([4, 3])
Step 0 completado: resultado = 7
...
```

**Servidores de Operação:**
```
Recebida operação: add([4, 3]) [Step: expr_1_step0]
Operação executada com sucesso: 7
```

## 🛑 **Encerramento**

Para parar todos os processos:

**Windows:**
```powershell
# Encontrar e matar processos Go
taskkill /F /IM grpc_add_server.exe
taskkill /F /IM grpc_sub_server.exe
taskkill /F /IM grpc_mult_server.exe
taskkill /F /IM grpc_div_server.exe
taskkill /F /IM grpc_dispatcher.exe
```

Ou simplesmente feche todos os terminais.

## 📝 **Notas Importantes**

1. **Ordem de inicialização importa**: Sempre inicie os servidores de operação antes do dispatcher
2. **Timeout padrão**: 30 segundos para cada operação
3. **Tratamento de erros**: O sistema trata divisão por zero e expressões inválidas
4. **Parsing**: Usa algoritmo Shunting Yard para converter infix → RPN
5. **Logs**: Todos os componentes geram logs detalhados para debug

## 📚 **Estrutura do Projeto**

```
ProjetoFinal/
├── cmd/
│   ├── grpc_client/          # Cliente CLI
│   ├── grpc_dispatcher/      # Dispatcher central
│   ├── grpc_add_server/      # Servidor Add
│   ├── grpc_sub_server/      # Servidor Subtract
│   ├── grpc_mult_server/     # Servidor Multiply
│   └── grpc_div_server/      # Servidor Divide
├── internal/
│   ├── core/                 # Parser e modelos
│   └── grpc/                 # Lógica de operações
├── proto/
│   └── calculator.proto      # Definições gRPC
├── bin/                      # Binários compilados
├── go.mod                    # Dependências
├── Makefile                  # Comandos de build
└── INSTRUCOES.md            # Este arquivo
```

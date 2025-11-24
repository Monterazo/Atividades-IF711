# 📋 **Instruções de Execução - Calculadora Distribuída gRPC**

Este documento fornece instruções completas para compilar e executar o sistema de calculadora distribuída em **Windows**, **Linux** e **macOS**.

---

## ⚡ **Quick Start (Resumo Rápido)**

### 🪟 Windows
```powershell
# 1. Instalar pré-requisitos
choco install golang protoc make  # ou instalar manualmente

# 2. Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Build do projeto
cd C:\Users\mikae\Documents\GitHub\Atividades-IF711\ProjetoFinal
go mod tidy
make proto   # ou comando protoc manual
make build   # ou comandos go build manuais

# 4. Executar (6 terminais diferentes)
.\bin\grpc_add_server.exe      # Terminal 1
.\bin\grpc_sub_server.exe      # Terminal 2
.\bin\grpc_mult_server.exe     # Terminal 3
.\bin\grpc_div_server.exe      # Terminal 4
.\bin\grpc_dispatcher.exe      # Terminal 5
.\bin\grpc_client.exe          # Terminal 6
```

### 🐧 Linux
```bash
# 1. Instalar pré-requisitos
sudo apt update && sudo apt install -y golang protobuf-compiler make

# 2. Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Build do projeto
cd ~/Atividades-IF711/ProjetoFinal
go mod tidy
make proto
make build

# 4. Executar (6 terminais diferentes)
./bin/grpc_add_server      # Terminal 1
./bin/grpc_sub_server      # Terminal 2
./bin/grpc_mult_server     # Terminal 3
./bin/grpc_div_server      # Terminal 4
./bin/grpc_dispatcher      # Terminal 5
./bin/grpc_client          # Terminal 6
```

### 🍎 macOS
```bash
# 1. Instalar pré-requisitos
brew install go protobuf

# 2. Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Build do projeto
cd ~/Atividades-IF711/ProjetoFinal
go mod tidy
make proto
make build

# 4. Executar (6 terminais diferentes)
./bin/grpc_add_server      # Terminal 1
./bin/grpc_sub_server      # Terminal 2
./bin/grpc_mult_server     # Terminal 3
./bin/grpc_div_server      # Terminal 4
./bin/grpc_dispatcher      # Terminal 5
./bin/grpc_client          # Terminal 6
```

---

## 🔧 **Pré-requisitos**

Antes de executar o projeto, certifique-se de ter instalado:

1. **Go 1.21+** - [Download](https://golang.org/dl/)
2. **Protocol Buffers Compiler (protoc)** - [Download](https://grpc.io/docs/protoc-installation/)
3. **Git** (para clonar o repositório)
4. **Make** (opcional, mas recomendado)

---

## 🖥️ **Instalação por Sistema Operacional**

### 🪟 **Windows**

#### 1. Instalar Go
```powershell
# Baixar e instalar Go de https://golang.org/dl/
# Verificar instalação
go version
```

#### 2. Instalar Protocol Buffers Compiler (protoc)
```powershell
# Opção 1: Usando Chocolatey (recomendado)
choco install protoc

# Opção 2: Manual
# 1. Baixar de https://github.com/protocolbuffers/protobuf/releases
# 2. Baixar o arquivo protoc-XX.X-win64.zip
# 3. Extrair para C:\protoc
# 4. Adicionar C:\protoc\bin ao PATH do sistema
```

#### 3. Instalar Make (opcional)
```powershell
# Usando Chocolatey
choco install make

# Ou usar comandos manuais (sem Makefile)
```

#### 4. Instalar plugins Go para protoc
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Verificar se está no PATH
# Deve estar em %USERPROFILE%\go\bin
```

---

### 🐧 **Linux (Ubuntu/Debian)**

#### 1. Instalar Go
```bash
# Opção 1: Usando gerenciador de pacotes
sudo apt update
sudo apt install golang-go

# Opção 2: Manual (versão mais recente)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verificar instalação
go version
```

#### 2. Instalar Protocol Buffers Compiler (protoc)
```bash
# Opção 1: Usando gerenciador de pacotes
sudo apt update
sudo apt install -y protobuf-compiler

# Opção 2: Manual (versão mais recente)
PB_VERSION="25.1"
wget https://github.com/protocolbuffers/protobuf/releases/download/v${PB_VERSION}/protoc-${PB_VERSION}-linux-x86_64.zip
sudo unzip protoc-${PB_VERSION}-linux-x86_64.zip -d /usr/local
sudo chmod +x /usr/local/bin/protoc

# Verificar instalação
protoc --version
```

#### 3. Instalar Make
```bash
sudo apt install make
```

#### 4. Instalar plugins Go para protoc
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicionar ao PATH se necessário
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

---

### 🍎 **macOS**

#### 1. Instalar Go
```bash
# Usando Homebrew (recomendado)
brew install go

# Verificar instalação
go version
```

#### 2. Instalar Protocol Buffers Compiler (protoc)
```bash
# Usando Homebrew
brew install protobuf

# Verificar instalação
protoc --version
```

#### 3. Instalar Make (geralmente já vem instalado)
```bash
# Verificar se já tem
make --version

# Se não tiver, instalar Xcode Command Line Tools
xcode-select --install
```

#### 4. Instalar plugins Go para protoc
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicionar ao PATH se necessário
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

---

## 📦 **Build do Projeto (Primeira Vez)**

Depois de instalar os pré-requisitos, siga os passos abaixo para compilar o projeto.

### 🪟 **Windows**

#### Passo 1: Navegar até o diretório do projeto
```powershell
cd C:\Users\mikae\Documents\GitHub\Atividades-IF711\ProjetoFinal
```

#### Passo 2: Baixar dependências do Go
```powershell
go mod download
go mod tidy
```

#### Passo 3: Gerar código a partir do arquivo .proto

**Opção A - Usando Make (se instalou):**
```powershell
make proto
```

**Opção B - Manual:**
```powershell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/calculator.proto
```

#### Passo 4: Compilar os binários

**Opção A - Usando Make:**
```powershell
make build
```

**Opção B - Manual:**
```powershell
# Criar diretório bin se não existir
if (-not (Test-Path "bin")) { New-Item -ItemType Directory -Path "bin" }

# Compilar cada componente
go build -o bin/grpc_add_server.exe cmd/grpc_add_server/main.go
go build -o bin/grpc_sub_server.exe cmd/grpc_sub_server/main.go
go build -o bin/grpc_mult_server.exe cmd/grpc_mult_server/main.go
go build -o bin/grpc_div_server.exe cmd/grpc_div_server/main.go
go build -o bin/grpc_dispatcher.exe cmd/grpc_dispatcher/main.go
go build -o bin/grpc_client.exe cmd/grpc_client/main.go
```

**Binários gerados em `bin/`:**
- `grpc_add_server.exe` - Servidor de adição
- `grpc_sub_server.exe` - Servidor de subtração
- `grpc_mult_server.exe` - Servidor de multiplicação
- `grpc_div_server.exe` - Servidor de divisão
- `grpc_dispatcher.exe` - Dispatcher central
- `grpc_client.exe` - Cliente CLI

---

### 🐧 **Linux**

#### Passo 1: Navegar até o diretório do projeto
```bash
cd ~/Atividades-IF711/ProjetoFinal
# ou onde você clonou o repositório
```

#### Passo 2: Baixar dependências do Go
```bash
go mod download
go mod tidy
```

#### Passo 3: Gerar código a partir do arquivo .proto

**Opção A - Usando Make:**
```bash
make proto
```

**Opção B - Manual:**
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/calculator.proto
```

#### Passo 4: Compilar os binários

**Opção A - Usando Make:**
```bash
make build
```

**Opção B - Manual:**
```bash
# Criar diretório bin se não existir
mkdir -p bin

# Compilar cada componente
go build -o bin/grpc_add_server cmd/grpc_add_server/main.go
go build -o bin/grpc_sub_server cmd/grpc_sub_server/main.go
go build -o bin/grpc_mult_server cmd/grpc_mult_server/main.go
go build -o bin/grpc_div_server cmd/grpc_div_server/main.go
go build -o bin/grpc_dispatcher cmd/grpc_dispatcher/main.go
go build -o bin/grpc_client cmd/grpc_client/main.go

# Dar permissão de execução
chmod +x bin/*
```

**Binários gerados em `bin/`:**
- `grpc_add_server` - Servidor de adição
- `grpc_sub_server` - Servidor de subtração
- `grpc_mult_server` - Servidor de multiplicação
- `grpc_div_server` - Servidor de divisão
- `grpc_dispatcher` - Dispatcher central
- `grpc_client` - Cliente CLI

---

### 🍎 **macOS**

#### Passo 1: Navegar até o diretório do projeto
```bash
cd ~/Atividades-IF711/ProjetoFinal
# ou onde você clonou o repositório
```

#### Passo 2: Baixar dependências do Go
```bash
go mod download
go mod tidy
```

#### Passo 3: Gerar código a partir do arquivo .proto

**Opção A - Usando Make:**
```bash
make proto
```

**Opção B - Manual:**
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/calculator.proto
```

#### Passo 4: Compilar os binários

**Opção A - Usando Make:**
```bash
make build
```

**Opção B - Manual:**
```bash
# Criar diretório bin se não existir
mkdir -p bin

# Compilar cada componente
go build -o bin/grpc_add_server cmd/grpc_add_server/main.go
go build -o bin/grpc_sub_server cmd/grpc_sub_server/main.go
go build -o bin/grpc_mult_server cmd/grpc_mult_server/main.go
go build -o bin/grpc_div_server cmd/grpc_div_server/main.go
go build -o bin/grpc_dispatcher cmd/grpc_dispatcher/main.go
go build -o bin/grpc_client cmd/grpc_client/main.go

# Dar permissão de execução
chmod +x bin/*
```

**Binários gerados em `bin/`:**
- `grpc_add_server` - Servidor de adição
- `grpc_sub_server` - Servidor de subtração
- `grpc_mult_server` - Servidor de multiplicação
- `grpc_div_server` - Servidor de divisão
- `grpc_dispatcher` - Dispatcher central
- `grpc_client` - Cliente CLI

---

## 🚀 **Execução do Sistema**

Após compilar, você precisa executar os componentes na ordem correta.

### 🪟 **Windows**

#### Opção 1: Execução Manual (Recomendado para debug/desenvolvimento)

Abra **6 terminais PowerShell diferentes** e execute em ordem:

**Terminal 1 - Servidor Add:**
```powershell
.\bin\grpc_add_server.exe
```

**Terminal 2 - Servidor Subtract:**
```powershell
.\bin\grpc_sub_server.exe
```

**Terminal 3 - Servidor Multiply:**
```powershell
.\bin\grpc_mult_server.exe
```

**Terminal 4 - Servidor Divide:**
```powershell
.\bin\grpc_div_server.exe
```

**Terminal 5 - Dispatcher:**
```powershell
.\bin\grpc_dispatcher.exe
```

**Terminal 6 - Cliente:**
```powershell
.\bin\grpc_client.exe
```

#### Opção 2: Execução Automatizada (Background)

**PowerShell:**
```powershell
# Iniciar todos os servidores em background
Start-Process -NoNewWindow .\bin\grpc_add_server.exe
Start-Process -NoNewWindow .\bin\grpc_sub_server.exe
Start-Process -NoNewWindow .\bin\grpc_mult_server.exe
Start-Process -NoNewWindow .\bin\grpc_div_server.exe

# Aguardar inicialização dos servidores
Start-Sleep -Seconds 2

# Iniciar dispatcher
Start-Process -NoNewWindow .\bin\grpc_dispatcher.exe

# Aguardar inicialização do dispatcher
Start-Sleep -Seconds 2

# Iniciar cliente (interativo)
.\bin\grpc_client.exe
```

#### Parar todos os processos (Windows):
```powershell
taskkill /F /IM grpc_add_server.exe
taskkill /F /IM grpc_sub_server.exe
taskkill /F /IM grpc_mult_server.exe
taskkill /F /IM grpc_div_server.exe
taskkill /F /IM grpc_dispatcher.exe
```

---

### 🐧 **Linux**

#### Opção 1: Execução Manual (Recomendado para debug/desenvolvimento)

Abra **6 terminais diferentes** e execute em ordem:

**Terminal 1 - Servidor Add:**
```bash
./bin/grpc_add_server
```

**Terminal 2 - Servidor Subtract:**
```bash
./bin/grpc_sub_server
```

**Terminal 3 - Servidor Multiply:**
```bash
./bin/grpc_mult_server
```

**Terminal 4 - Servidor Divide:**
```bash
./bin/grpc_div_server
```

**Terminal 5 - Dispatcher:**
```bash
./bin/grpc_dispatcher
```

**Terminal 6 - Cliente:**
```bash
./bin/grpc_client
```

#### Opção 2: Execução Automatizada (Background)

```bash
# Iniciar todos os servidores em background
./bin/grpc_add_server > /tmp/grpc_add.log 2>&1 &
./bin/grpc_sub_server > /tmp/grpc_sub.log 2>&1 &
./bin/grpc_mult_server > /tmp/grpc_mult.log 2>&1 &
./bin/grpc_div_server > /tmp/grpc_div.log 2>&1 &

# Aguardar inicialização dos servidores
sleep 2

# Iniciar dispatcher em background
./bin/grpc_dispatcher > /tmp/grpc_dispatcher.log 2>&1 &

# Aguardar inicialização do dispatcher
sleep 2

# Iniciar cliente (interativo)
./bin/grpc_client
```

#### Opção 3: Usando Make (se disponível):
```bash
make run-all
```

#### Parar todos os processos (Linux):
```bash
pkill -f grpc_add_server
pkill -f grpc_sub_server
pkill -f grpc_mult_server
pkill -f grpc_div_server
pkill -f grpc_dispatcher
```

---

### 🍎 **macOS**

#### Opção 1: Execução Manual (Recomendado para debug/desenvolvimento)

Abra **6 terminais diferentes** e execute em ordem:

**Terminal 1 - Servidor Add:**
```bash
./bin/grpc_add_server
```

**Terminal 2 - Servidor Subtract:**
```bash
./bin/grpc_sub_server
```

**Terminal 3 - Servidor Multiply:**
```bash
./bin/grpc_mult_server
```

**Terminal 4 - Servidor Divide:**
```bash
./bin/grpc_div_server
```

**Terminal 5 - Dispatcher:**
```bash
./bin/grpc_dispatcher
```

**Terminal 6 - Cliente:**
```bash
./bin/grpc_client
```

#### Opção 2: Execução Automatizada (Background)

```bash
# Iniciar todos os servidores em background
./bin/grpc_add_server > /tmp/grpc_add.log 2>&1 &
./bin/grpc_sub_server > /tmp/grpc_sub.log 2>&1 &
./bin/grpc_mult_server > /tmp/grpc_mult.log 2>&1 &
./bin/grpc_div_server > /tmp/grpc_div.log 2>&1 &

# Aguardar inicialização dos servidores
sleep 2

# Iniciar dispatcher em background
./bin/grpc_dispatcher > /tmp/grpc_dispatcher.log 2>&1 &

# Aguardar inicialização do dispatcher
sleep 2

# Iniciar cliente (interativo)
./bin/grpc_client
```

#### Opção 3: Usando Make:
```bash
make run-all
```

#### Parar todos os processos (macOS):
```bash
pkill -f grpc_add_server
pkill -f grpc_sub_server
pkill -f grpc_mult_server
pkill -f grpc_div_server
pkill -f grpc_dispatcher
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

---

## 🐛 **Troubleshooting (Solução de Problemas)**

### ❌ Problemas Comuns (Todas as Plataformas)

#### Erro: "Falha ao conectar ao dispatcher"
```
Possíveis causas:
- Dispatcher não está rodando
- Porta 50051 está sendo usada por outro processo
```

**Solução:**
```bash
# Verificar se a porta está em uso

# Windows:
netstat -ano | findstr :50051

# Linux/macOS:
lsof -i :50051

# Matar processo que está usando a porta (se necessário)
# Windows:
taskkill /PID <PID> /F

# Linux/macOS:
kill -9 <PID>
```

#### Erro: "Falha ao conectar aos servidores"
```
Possíveis causas:
- Servidores de operação não estão rodando
- Portas 50052-50055 estão sendo usadas
```

**Solução:**
- Certifique-se de que todos os 4 servidores foram iniciados
- Verifique se as portas estão livres (use comandos acima para cada porta)
- Inicie os servidores ANTES do dispatcher

#### Erro: "protoc: command not found" ou "protoc-gen-go: command not found"

**Windows:**
```powershell
# Verificar se protoc está instalado
protoc --version

# Verificar se está no PATH
$env:PATH -split ';' | Select-String protoc

# Adicionar ao PATH se necessário
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\protoc\bin", "User")

# Verificar plugins Go
go env GOPATH
# Adicionar %GOPATH%\bin ao PATH se necessário
```

**Linux/macOS:**
```bash
# Verificar instalação
which protoc
protoc --version

# Verificar plugins Go
which protoc-gen-go
which protoc-gen-go-grpc

# Adicionar ao PATH se necessário
export PATH=$PATH:$(go env GOPATH)/bin
# Adicionar ao ~/.bashrc ou ~/.zshrc para permanência
```

#### Erro ao compilar .proto
```bash
# Reinstalar plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Verificar se foram instalados
# Windows:
where protoc-gen-go

# Linux/macOS:
which protoc-gen-go
```

### 🪟 Problemas Específicos do Windows

#### Erro: "cannot execute binary file"
- Certifique-se de baixar a versão **win64** do protoc
- Não use arquivos compilados para Linux/macOS

#### PowerShell não reconhece comandos
```powershell
# Executar PowerShell como Administrador
# Permitir execução de scripts
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### Make não encontrado
```powershell
# Instalar Make via Chocolatey
choco install make

# Ou executar comandos manualmente sem Makefile
```

### 🐧 Problemas Específicos do Linux

#### Erro: "permission denied" ao executar binários
```bash
# Dar permissão de execução
chmod +x bin/*
```

#### Erro: "protoc: error while loading shared libraries"
```bash
# Instalar dependências
sudo apt install -y libprotobuf-dev

# Ou recompilar protoc manualmente
```

#### Go não está no PATH
```bash
# Adicionar ao ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### 🍎 Problemas Específicos do macOS

#### Erro: "command not found" mesmo após instalar via brew
```bash
# Recarregar shell
source ~/.zshrc

# Verificar se Homebrew está no PATH
echo $PATH | grep homebrew

# Adicionar Homebrew ao PATH se necessário
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zshrc
source ~/.zshrc
```

#### Erro: "xcrun: error: invalid active developer path"
```bash
# Instalar Xcode Command Line Tools
xcode-select --install
```

#### Protoc instalado mas não funciona
```bash
# Reinstalar via Homebrew
brew uninstall protobuf
brew install protobuf

# Verificar instalação
protoc --version
```

### 🔧 Problemas de Dependências Go

#### Erro: "package X is not in GOROOT"
```bash
# Atualizar dependências
go mod download
go mod tidy

# Limpar cache se necessário
go clean -modcache
```

#### Erro: "go: module X: Get ... connection refused"
```bash
# Configurar proxy Go (útil em algumas regiões)
go env -w GOPROXY=https://goproxy.io,direct

# Ou usar proxy alternativo
go env -w GOPROXY=https://proxy.golang.org,direct
```

### 📊 Verificar Logs em Caso de Erro

**Windows:**
```powershell
# Logs aparecem no terminal onde o processo está rodando
# Para executar com log em arquivo:
.\bin\grpc_dispatcher.exe > dispatcher.log 2>&1
```

**Linux/macOS:**
```bash
# Verificar logs quando executado em background
tail -f /tmp/grpc_dispatcher.log
tail -f /tmp/grpc_add.log

# Verificar todos os logs
tail -f /tmp/grpc_*.log
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

---

## 🛑 **Encerramento do Sistema**

### 🪟 **Windows**

#### Método 1: Fechar terminais
Simplesmente feche todos os 6 terminais PowerShell que você abriu.

#### Método 2: Matar processos via comando
```powershell
# Matar todos os processos gRPC
taskkill /F /IM grpc_add_server.exe
taskkill /F /IM grpc_sub_server.exe
taskkill /F /IM grpc_mult_server.exe
taskkill /F /IM grpc_div_server.exe
taskkill /F /IM grpc_dispatcher.exe
taskkill /F /IM grpc_client.exe
```

#### Método 3: Verificar e limpar processos órfãos
```powershell
# Listar todos os processos gRPC
tasklist | findstr grpc_

# Matar por PID se necessário
taskkill /F /PID <PID>
```

---

### 🐧 **Linux**

#### Método 1: Ctrl+C nos terminais
Pressione `Ctrl+C` em cada terminal onde os processos estão rodando.

#### Método 2: Matar processos via comando
```bash
# Matar todos os processos gRPC
pkill -f grpc_add_server
pkill -f grpc_sub_server
pkill -f grpc_mult_server
pkill -f grpc_div_server
pkill -f grpc_dispatcher
pkill -f grpc_client
```

#### Método 3: Verificar e limpar processos órfãos
```bash
# Listar todos os processos gRPC
ps aux | grep grpc_

# Matar por PID se necessário
kill -9 <PID>

# Ou matar todos de uma vez
killall grpc_add_server grpc_sub_server grpc_mult_server grpc_div_server grpc_dispatcher grpc_client
```

---

### 🍎 **macOS**

#### Método 1: Ctrl+C nos terminais
Pressione `Cmd+C` ou `Ctrl+C` em cada terminal onde os processos estão rodando.

#### Método 2: Matar processos via comando
```bash
# Matar todos os processos gRPC
pkill -f grpc_add_server
pkill -f grpc_sub_server
pkill -f grpc_mult_server
pkill -f grpc_div_server
pkill -f grpc_dispatcher
pkill -f grpc_client
```

#### Método 3: Verificar e limpar processos órfãos
```bash
# Listar todos os processos gRPC
ps aux | grep grpc_

# Matar por PID se necessário
kill -9 <PID>

# Ou matar todos de uma vez
killall grpc_add_server grpc_sub_server grpc_mult_server grpc_div_server grpc_dispatcher grpc_client
```

---

## 📝 **Notas Importantes**

1. **Ordem de inicialização importa**: Sempre inicie os servidores de operação ANTES do dispatcher
2. **Timeout padrão**: 30 segundos para cada operação
3. **Tratamento de erros**: O sistema trata divisão por zero e expressões inválidas
4. **Parsing**: Usa algoritmo Shunting Yard para converter infix → RPN
5. **Logs**: Todos os componentes geram logs detalhados para debug
6. **Compatibilidade**: Testado em Windows 10/11, Ubuntu 20.04+, macOS 11+

---

## ✅ **Checklist de Verificação Rápida**

Antes de reportar problemas, verifique:

### Pré-requisitos Instalados?
- [ ] Go 1.21+ instalado (`go version`)
- [ ] Protoc instalado (`protoc --version`)
- [ ] Plugins Go instalados (`which protoc-gen-go` ou `where protoc-gen-go`)
- [ ] Variáveis de ambiente PATH configuradas corretamente

### Build Realizado?
- [ ] `go mod download` executado sem erros
- [ ] `go mod tidy` executado sem erros
- [ ] Código .proto gerado (pasta `proto/` contém arquivos `.pb.go`)
- [ ] Binários compilados (pasta `bin/` contém executáveis)

### Execução Correta?
- [ ] Todos os 4 servidores de operação rodando
- [ ] Dispatcher rodando (iniciado DEPOIS dos servidores)
- [ ] Portas 50051-50055 não estão sendo usadas por outros processos
- [ ] Cliente consegue se conectar ao dispatcher

### Em Caso de Erro:
- [ ] Logs dos servidores verificados
- [ ] Processos órfãos eliminados
- [ ] Portas liberadas
- [ ] Dependências atualizadas (`go mod tidy`)

---

## 📑 **Índice de Navegação Rápida**

| Seção | Link | Descrição |
|-------|------|-----------|
| **Quick Start** | [↑ Ir para seção](#-quick-start-resumo-rápido) | Comandos resumidos por plataforma |
| **Instalação Windows** | [↑ Ir para seção](#-windows) | Instalação completa no Windows |
| **Instalação Linux** | [↑ Ir para seção](#-linux-ubuntudebian) | Instalação completa no Linux |
| **Instalação macOS** | [↑ Ir para seção](#-macos) | Instalação completa no macOS |
| **Build** | [↑ Ir para seção](#-build-do-projeto-primeira-vez) | Compilar o projeto |
| **Execução** | [↑ Ir para seção](#-execução-do-sistema) | Executar os componentes |
| **Usando o Cliente** | [↑ Ir para seção](#-usando-o-cliente) | Como usar a interface CLI |
| **Troubleshooting** | [↑ Ir para seção](#-troubleshooting-solução-de-problemas) | Solução de problemas comuns |
| **Testes** | [↑ Ir para seção](#-testes) | Exemplos de testes |
| **Encerramento** | [↑ Ir para seção](#-encerramento-do-sistema) | Como parar o sistema |

---

## 📞 **Suporte e Documentação Adicional**

- **Repositório**: [GitHub - Atividades IF711](https://github.com/...)
- **Issues**: Reporte problemas no GitHub Issues
- **Documentação gRPC**: https://grpc.io/docs/
- **Documentação Go**: https://go.dev/doc/
- **Protocol Buffers**: https://protobuf.dev/

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

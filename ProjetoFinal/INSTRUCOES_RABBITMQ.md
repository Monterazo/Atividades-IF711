# 📋 **Instruções de Execução - Calculadora Distribuída RabbitMQ**

Este documento fornece instruções completas para compilar e executar o sistema de calculadora distribuída baseado em **RabbitMQ** em **Windows**, **Linux** e **macOS**.

---

## ⚡ **Quick Start (Resumo Rápido)**

### 🪟 Windows
```powershell
# 1. Instalar RabbitMQ
choco install rabbitmq

# 2. Iniciar RabbitMQ
rabbitmq-server

# 3. Build do projeto
cd C:\Users\mikae\Documents\GitHub\Atividades-IF711\ProjetoFinal
go mod tidy
make build-rabbitmq

# 4. Executar (6 terminais diferentes)
.\bin\rabbitmq_add_server.exe      # Terminal 1
.\bin\rabbitmq_sub_server.exe      # Terminal 2
.\bin\rabbitmq_mult_server.exe     # Terminal 3
.\bin\rabbitmq_div_server.exe      # Terminal 4
.\bin\rabbitmq_dispatcher.exe      # Terminal 5
.\bin\rabbitmq_client.exe          # Terminal 6
```

### 🐧 Linux
```bash
# 1. Instalar RabbitMQ
sudo apt update && sudo apt install -y rabbitmq-server

# 2. Iniciar RabbitMQ
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server

# 3. Build do projeto
cd ~/Atividades-IF711/ProjetoFinal
go mod tidy
make build-rabbitmq

# 4. Executar (6 terminais diferentes)
./bin/rabbitmq_add_server      # Terminal 1
./bin/rabbitmq_sub_server      # Terminal 2
./bin/rabbitmq_mult_server     # Terminal 3
./bin/rabbitmq_div_server      # Terminal 4
./bin/rabbitmq_dispatcher      # Terminal 5
./bin/rabbitmq_client          # Terminal 6
```

### 🍎 macOS
```bash
# 1. Instalar RabbitMQ
brew install rabbitmq

# 2. Iniciar RabbitMQ
brew services start rabbitmq

# 3. Build do projeto
cd ~/Atividades-IF711/ProjetoFinal
go mod tidy
make build-rabbitmq

# 4. Executar (6 terminais diferentes)
./bin/rabbitmq_add_server      # Terminal 1
./bin/rabbitmq_sub_server      # Terminal 2
./bin/rabbitmq_mult_server     # Terminal 3
./bin/rabbitmq_div_server      # Terminal 4
./bin/rabbitmq_dispatcher      # Terminal 5
./bin/rabbitmq_client          # Terminal 6
```

---

## 🔧 **Pré-requisitos**

Antes de executar o projeto, certifique-se de ter instalado:

1. **Go 1.21+** - [Download](https://golang.org/dl/)
2. **RabbitMQ Server** - [Download](https://www.rabbitmq.com/download.html)
3. **Git** (para clonar o repositório)
4. **Make** (opcional, mas recomendado)

---

## 🐰 **Instalação do RabbitMQ por Sistema Operacional**

### 🪟 **Windows**

#### Método 1: Usando Chocolatey (Recomendado)
```powershell
# Instalar Erlang (pré-requisito do RabbitMQ)
choco install erlang

# Instalar RabbitMQ
choco install rabbitmq

# Verificar instalação
rabbitmq-diagnostics status
```

#### Método 2: Manual
1. Baixar Erlang de https://www.erlang.org/downloads
2. Baixar RabbitMQ de https://www.rabbitmq.com/download.html
3. Instalar Erlang primeiro, depois RabbitMQ
4. Adicionar ao PATH: `C:\Program Files\RabbitMQ Server\rabbitmq_server-X.X.X\sbin`

#### Iniciar RabbitMQ
```powershell
# Iniciar servidor
rabbitmq-server

# Ou como serviço Windows
rabbitmq-service start

# Verificar status
rabbitmq-diagnostics status
```

#### Habilitar Management UI (Opcional)
```powershell
rabbitmq-plugins enable rabbitmq_management

# Acessar em: http://localhost:15672
# Usuário: guest
# Senha: guest
```

---

### 🐧 **Linux (Ubuntu/Debian)**

#### Instalação
```bash
# Atualizar repositórios
sudo apt update

# Instalar RabbitMQ
sudo apt install -y rabbitmq-server

# Verificar status
sudo systemctl status rabbitmq-server
```

#### Iniciar RabbitMQ
```bash
# Iniciar servidor
sudo systemctl start rabbitmq-server

# Habilitar inicialização automática
sudo systemctl enable rabbitmq-server

# Verificar status
rabbitmq-diagnostics status
```

#### Habilitar Management UI (Opcional)
```bash
sudo rabbitmq-plugins enable rabbitmq_management

# Criar usuário admin (opcional)
sudo rabbitmqctl add_user admin admin
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"

# Acessar em: http://localhost:15672
```

---

### 🍎 **macOS**

#### Instalação usando Homebrew
```bash
# Instalar RabbitMQ
brew install rabbitmq

# Adicionar ao PATH (adicionar ao ~/.zshrc ou ~/.bash_profile)
export PATH=$PATH:/opt/homebrew/opt/rabbitmq/sbin

# Recarregar configuração
source ~/.zshrc
```

#### Iniciar RabbitMQ
```bash
# Iniciar como serviço
brew services start rabbitmq

# Ou executar manualmente
rabbitmq-server

# Verificar status
rabbitmq-diagnostics status
```

#### Habilitar Management UI (Opcional)
```bash
rabbitmq-plugins enable rabbitmq_management

# Acessar em: http://localhost:15672
# Usuário: guest
# Senha: guest
```

---

## 📦 **Build do Projeto**

### Passo 1: Navegar até o diretório do projeto
```bash
cd ProjetoFinal
```

### Passo 2: Baixar dependências do Go
```bash
go mod download
go mod tidy
```

### Passo 3: Compilar os binários RabbitMQ

**Usando Make:**
```bash
make build-rabbitmq
```

**Manual (Windows PowerShell):**
```powershell
# Criar diretório bin se não existir
if (-not (Test-Path "bin")) { New-Item -ItemType Directory -Path "bin" }

# Compilar cada componente
go build -o bin/rabbitmq_add_server.exe cmd/rabbitmq_add_server/main.go
go build -o bin/rabbitmq_sub_server.exe cmd/rabbitmq_sub_server/main.go
go build -o bin/rabbitmq_mult_server.exe cmd/rabbitmq_mult_server/main.go
go build -o bin/rabbitmq_div_server.exe cmd/rabbitmq_div_server/main.go
go build -o bin/rabbitmq_dispatcher.exe cmd/rabbitmq_dispatcher/main.go
go build -o bin/rabbitmq_client.exe cmd/rabbitmq_client/main.go
```

**Manual (Linux/macOS):**
```bash
# Criar diretório bin se não existir
mkdir -p bin

# Compilar cada componente
go build -o bin/rabbitmq_add_server cmd/rabbitmq_add_server/main.go
go build -o bin/rabbitmq_sub_server cmd/rabbitmq_sub_server/main.go
go build -o bin/rabbitmq_mult_server cmd/rabbitmq_mult_server/main.go
go build -o bin/rabbitmq_div_server cmd/rabbitmq_div_server/main.go
go build -o bin/rabbitmq_dispatcher cmd/rabbitmq_dispatcher/main.go
go build -o bin/rabbitmq_client cmd/rabbitmq_client/main.go

# Dar permissão de execução
chmod +x bin/rabbitmq_*
```

**Binários gerados em `bin/`:**
- `rabbitmq_add_server` - Servidor de adição
- `rabbitmq_sub_server` - Servidor de subtração
- `rabbitmq_mult_server` - Servidor de multiplicação
- `rabbitmq_div_server` - Servidor de divisão
- `rabbitmq_dispatcher` - Dispatcher central
- `rabbitmq_client` - Cliente CLI

---

## 🚀 **Execução do Sistema**

**IMPORTANTE:** Certifique-se de que o RabbitMQ está rodando antes de iniciar os componentes!

### 🪟 **Windows**

#### Opção 1: Execução Manual (Recomendado)

Abra **6 terminais PowerShell diferentes** e execute em ordem:

**Terminal 1 - Servidor Add:**
```powershell
.\bin\rabbitmq_add_server.exe
```

**Terminal 2 - Servidor Subtract:**
```powershell
.\bin\rabbitmq_sub_server.exe
```

**Terminal 3 - Servidor Multiply:**
```powershell
.\bin\rabbitmq_mult_server.exe
```

**Terminal 4 - Servidor Divide:**
```powershell
.\bin\rabbitmq_div_server.exe
```

**Terminal 5 - Dispatcher:**
```powershell
.\bin\rabbitmq_dispatcher.exe
```

**Terminal 6 - Cliente:**
```powershell
.\bin\rabbitmq_client.exe
```

#### Opção 2: Usando Make
```powershell
make run-all-rabbitmq
```

---

### 🐧 **Linux**

#### Execução Manual

Abra **6 terminais diferentes** e execute:

**Terminal 1 - Servidor Add:**
```bash
./bin/rabbitmq_add_server
```

**Terminal 2 - Servidor Subtract:**
```bash
./bin/rabbitmq_sub_server
```

**Terminal 3 - Servidor Multiply:**
```bash
./bin/rabbitmq_mult_server
```

**Terminal 4 - Servidor Divide:**
```bash
./bin/rabbitmq_div_server
```

**Terminal 5 - Dispatcher:**
```bash
./bin/rabbitmq_dispatcher
```

**Terminal 6 - Cliente:**
```bash
./bin/rabbitmq_client
```

#### Usando Make
```bash
make run-all-rabbitmq
```

---

### 🍎 **macOS**

Siga os mesmos passos do Linux acima.

---

## 💻 **Usando o Cliente**

O cliente RabbitMQ funciona da mesma forma que o cliente gRPC:

```
Cliente Calculadora RabbitMQ
========================
Conectado ao RabbitMQ em amqp://guest:guest@localhost:5672/

Digite uma expressão matemática (ou 'sair' para encerrar):
Exemplos: ((4+3)*2)/5, 10+20*3, (15-5)/2

>
```

### Exemplos de Expressões:

```
> ((4+3)*2)/5
✅ Resultado: ((4+3)*2)/5 = 2.800000
⏱️  Tempo de execução: 25ms

> 10+20*3
✅ Resultado: 10+20*3 = 70.000000
⏱️  Tempo de execução: 18ms

> (15-5)/2
✅ Resultado: (15-5)/2 = 5.000000
⏱️  Tempo de execução: 15ms

> 10/0
❌ Erro: [DIV_BY_ZERO] divisão por zero
```

---

## 🏗️ **Arquitetura RabbitMQ**

```
Cliente
    │
    ▼ (publica em calculator.requests)
RabbitMQ Broker
    │
    ▼ (consome calculator.requests)
Dispatcher
    │
    ├──> (publica em operations.add) ──> AddServer
    ├──> (publica em operations.subtract) ──> SubServer
    ├──> (publica em operations.multiply) ──> MultServer
    └──> (publica em operations.divide) ──> DivServer
         │
         └──> (todos publicam em operations.results)
              │
              ▼ (dispatcher consome operations.results)
         Dispatcher
              │
              ▼ (publica em calculator.responses)
         RabbitMQ Broker
              │
              ▼ (cliente consome calculator.responses)
         Cliente
```

### Filas Utilizadas:
- **calculator.requests** - Requisições de expressões
- **calculator.responses** - Respostas para o cliente
- **operations.add** - Operações de adição
- **operations.subtract** - Operações de subtração
- **operations.multiply** - Operações de multiplicação
- **operations.divide** - Operações de divisão
- **operations.results** - Resultados das operações

---

## 🔍 **Fluxo de Execução**

1. **Cliente** publica expressão em `calculator.requests`
2. **Dispatcher** consome mensagem, faz parsing (Shunting Yard → RPN)
3. **Dispatcher** decompõe em steps e publica em filas específicas (`operations.add`, etc.)
4. Cada **Servidor** consome sua fila, executa operação e publica resultado em `operations.results`
5. **Dispatcher** consome resultados, agrupa e publica resposta final em `calculator.responses`
6. **Cliente** consome resposta e exibe resultado

---

## 🐛 **Troubleshooting**

### RabbitMQ não está rodando
```bash
# Linux/macOS
sudo systemctl status rabbitmq-server
sudo systemctl start rabbitmq-server

# macOS (Homebrew)
brew services list
brew services start rabbitmq

# Windows
rabbitmq-service status
rabbitmq-service start
```

### Erro: "Falha ao conectar ao RabbitMQ"
```bash
# Verificar se RabbitMQ está escutando na porta 5672
# Linux/macOS:
sudo lsof -i :5672

# Windows:
netstat -ano | findstr :5672

# Verificar logs do RabbitMQ
# Linux: /var/log/rabbitmq/
# macOS: /opt/homebrew/var/log/rabbitmq/
# Windows: C:\Users\<user>\AppData\Roaming\RabbitMQ\log\
```

### Limpar filas do RabbitMQ
```bash
# Via CLI
rabbitmqctl purge_queue calculator.requests
rabbitmqctl purge_queue calculator.responses
rabbitmqctl purge_queue operations.results

# Ou deletar e recriar todas as filas
rabbitmqctl list_queues
rabbitmqctl delete_queue <nome_da_fila>
```

### Monitorar filas via Management UI
1. Acesse http://localhost:15672
2. Login: guest / guest
3. Vá para aba "Queues"
4. Visualize mensagens, consumidores e taxas

---

## 🛑 **Encerramento do Sistema**

### Parar componentes
```bash
# Pressionar Ctrl+C em cada terminal

# Ou matar processos (Linux/macOS)
pkill -f rabbitmq_add_server
pkill -f rabbitmq_sub_server
pkill -f rabbitmq_mult_server
pkill -f rabbitmq_div_server
pkill -f rabbitmq_dispatcher
pkill -f rabbitmq_client

# Windows
taskkill /F /IM rabbitmq_add_server.exe
taskkill /F /IM rabbitmq_sub_server.exe
taskkill /F /IM rabbitmq_mult_server.exe
taskkill /F /IM rabbitmq_div_server.exe
taskkill /F /IM rabbitmq_dispatcher.exe
taskkill /F /IM rabbitmq_client.exe
```

### Parar RabbitMQ
```bash
# Linux
sudo systemctl stop rabbitmq-server

# macOS
brew services stop rabbitmq

# Windows
rabbitmq-service stop
```

---

## 📊 **Comparação: RabbitMQ vs gRPC**

| Aspecto | RabbitMQ (MOM) | gRPC (RPC) |
|---------|----------------|------------|
| **Paradigma** | Assíncrono, baseado em mensagens | Síncrono, chamadas diretas |
| **Acoplamento** | Baixo (desacoplado via broker) | Alto (cliente conhece servidor) |
| **Latência** | Maior (overhead do broker) | Menor (conexão direta) |
| **Confiabilidade** | Alta (mensagens persistentes) | Média (depende da rede) |
| **Escalabilidade** | Alta (fácil adicionar consumidores) | Média (requer load balancer) |
| **Complexidade** | Média (requer broker) | Baixa (ponto-a-ponto) |
| **Falhas** | Tolera falhas temporárias | Falha imediata se servidor offline |

---

## 📝 **Notas Importantes**

1. **RabbitMQ deve estar rodando** antes de iniciar qualquer componente
2. **Ordem de inicialização**: Servidores → Dispatcher → Cliente
3. **Filas são duráveis**: Mensagens sobrevivem a reinicializações do RabbitMQ
4. **Timeout padrão**: 30 segundos para cada operação
5. **Formato de mensagens**: JSON
6. **Porta padrão do RabbitMQ**: 5672 (AMQP)
7. **Management UI**: Porta 15672 (HTTP)

---

## 🧪 **Testes**

Use os mesmos testes da versão gRPC:

### Operações Básicas:
```
> 5+3          # Esperado: 8
> 10-4         # Esperado: 6
> 6*7          # Esperado: 42
> 15/3         # Esperado: 5
```

### Expressões Complexas:
```
> ((4+3)*2)/5       # Esperado: 2.8
> 10+20*3           # Esperado: 70
> (15-5)/2          # Esperado: 5
```

### Teste de Erros:
```
> 10/0              # Esperado: Erro DIV_BY_ZERO
```

---

## 📚 **Recursos Adicionais**

- **Documentação RabbitMQ**: https://www.rabbitmq.com/documentation.html
- **Tutorial Go AMQP**: https://www.rabbitmq.com/tutorials/tutorial-one-go.html
- **Management Plugin**: https://www.rabbitmq.com/management.html
- **Monitoring**: https://www.rabbitmq.com/monitoring.html

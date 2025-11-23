# 🔧 **Guia de Configuração do Ambiente**

## ⚙️ **Instalação do Protocol Buffers Compiler (protoc)**

### Windows

#### Opção 1: Download Manual
1. Acesse: https://github.com/protocolbuffers/protobuf/releases
2. Baixe a versão mais recente para Windows (ex: `protoc-25.1-win64.zip`)
3. Extraia o arquivo ZIP
4. Adicione a pasta `bin` ao PATH do sistema:
   - Pressione `Win + R`
   - Digite `sysdm.cpl` e pressione Enter
   - Vá para "Variáveis de Ambiente"
   - Edite a variável `Path`
   - Adicione o caminho: `C:\protoc\bin` (ajuste conforme necessário)
5. Verifique a instalação:
   ```bash
   protoc --version
   ```

#### Opção 2: Usando Chocolatey
```bash
choco install protoc
```

#### Opção 3: Usando Scoop
```bash
scoop install protobuf
```

### Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y protobuf-compiler

# Fedora
sudo dnf install protobuf-compiler

# Arch Linux
sudo pacman -S protobuf
```

### macOS
```bash
brew install protobuf
```

## 📦 **Instalação dos Plugins Go**

Após instalar o `protoc`, instale os plugins Go:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Verifique se os plugins estão no PATH:
```bash
# Windows
echo %GOPATH%\bin

# Linux/macOS
echo $GOPATH/bin
```

Se não estiverem no PATH, adicione:

**Windows (PowerShell):**
```powershell
$env:PATH += ";$env:GOPATH\bin"
```

**Linux/macOS:**
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## 🚀 **Configuração do Projeto**

### 1. Baixar Dependências
```bash
cd C:\Users\mikae\Documents\GitHub\Atividades-IF711\ProjetoFinal
go mod download
go mod tidy
```

### 2. Gerar Código a partir do .proto

**Sem protoc instalado** (Alternativa):
Se você não conseguir instalar o protoc, pode usar uma imagem Docker:

```bash
docker run --rm -v ${PWD}:/workspace \
  namely/protoc-all \
  -f proto/calculator.proto \
  -l go \
  -o .
```

**Com protoc instalado**:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/calculator.proto
```

Isso irá gerar:
- `proto/calculator.pb.go` - Structs e tipos
- `proto/calculator_grpc.pb.go` - Interfaces de serviço

### 3. Compilar o Projeto
```bash
go build -o bin/grpc_add_server.exe ./cmd/grpc_add_server
go build -o bin/grpc_sub_server.exe ./cmd/grpc_sub_server
go build -o bin/grpc_mult_server.exe ./cmd/grpc_mult_server
go build -o bin/grpc_div_server.exe ./cmd/grpc_div_server
go build -o bin/grpc_dispatcher.exe ./cmd/grpc_dispatcher
go build -o bin/grpc_client.exe ./cmd/grpc_client
```

Ou use o Makefile:
```bash
make build
```

## ✅ **Verificação**

Para verificar se tudo está configurado corretamente:

```bash
# Verificar Go
go version

# Verificar protoc
protoc --version

# Verificar plugins
protoc-gen-go --version
protoc-gen-go-grpc --version

# Verificar dependências
go mod verify
```

## 🐛 **Problemas Comuns**

### "protoc: command not found"
- Certifique-se de que o protoc está no PATH
- Reinicie o terminal após adicionar ao PATH
- Use a alternativa com Docker

### "protoc-gen-go: program not found"
```bash
# Reinstale os plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Adicione GOPATH/bin ao PATH
export PATH=$PATH:$(go env GOPATH)/bin  # Linux/macOS
$env:PATH += ";$env:GOPATH\bin"         # Windows PowerShell
```

### Erro de importação no Go
```bash
# Limpe o cache e baixe novamente
go clean -modcache
go mod download
go mod tidy
```

### Porta já em uso
```bash
# Windows - Encontrar processo na porta 50051
netstat -ano | findstr :50051
taskkill /PID <PID> /F

# Linux/macOS
lsof -i :50051
kill -9 <PID>
```

## 📋 **Checklist de Configuração**

- [ ] Go 1.21+ instalado
- [ ] protoc instalado e no PATH
- [ ] protoc-gen-go instalado
- [ ] protoc-gen-go-grpc instalado
- [ ] Dependências do Go baixadas (`go mod download`)
- [ ] Código gerado a partir do .proto
- [ ] Projeto compilado com sucesso
- [ ] Portas 50051-50055 livres

## 🎯 **Próximos Passos**

Após concluir a configuração:
1. Consulte `INSTRUCOES.md` para instruções de execução
2. Leia `README.md` para entender a arquitetura
3. Execute os servidores e teste o sistema

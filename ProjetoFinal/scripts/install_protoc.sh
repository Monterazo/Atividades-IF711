#!/bin/bash
# Script de instalação do protoc para Git Bash (Windows)
# Não requer privilégios de administrador

set -e

PROTOC_VERSION="25.1"
INSTALL_DIR="$HOME/.local/protoc"
BIN_DIR="$INSTALL_DIR/bin"
TEMP_ZIP="/tmp/protoc.zip"

echo "======================================"
echo "Instalando protoc v$PROTOC_VERSION"
echo "======================================"

# Criar diretório de instalação
echo "Criando diretório de instalação em $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"

# Baixar protoc
echo "Baixando protoc v$PROTOC_VERSION..."
PROTOC_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-win64.zip"
curl -L "$PROTOC_URL" -o "$TEMP_ZIP"

# Extrair
echo "Extraindo arquivos..."
unzip -o "$TEMP_ZIP" -d "$INSTALL_DIR"

# Limpar
rm -f "$TEMP_ZIP"

# Adicionar ao PATH permanentemente no .bashrc
BASHRC="$HOME/.bashrc"
if ! grep -q "protoc/bin" "$BASHRC" 2>/dev/null; then
    echo "" >> "$BASHRC"
    echo "# Protoc PATH" >> "$BASHRC"
    echo "export PATH=\"\$HOME/.local/protoc/bin:\$PATH\"" >> "$BASHRC"
    echo "✅ Adicionado ao .bashrc"
else
    echo "✅ Já está no .bashrc"
fi

# Adicionar ao PATH da sessão atual
export PATH="$BIN_DIR:$PATH"

# Verificar instalação
echo ""
echo "Verificando instalação..."
"$BIN_DIR/protoc.exe" --version

echo ""
echo "======================================"
echo "✅ protoc instalado com sucesso!"
echo "======================================"
echo ""
echo "Instalando plugins Go..."
echo ""

# Instalar plugins Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo ""
echo "======================================"
echo "✅ Instalação completa!"
echo "======================================"
echo ""
echo "Próximos passos:"
echo "1. Execute: source ~/.bashrc"
echo "2. Execute: make proto"
echo ""

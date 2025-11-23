# Script de build para Windows
# Execute: powershell -ExecutionPolicy Bypass -File scripts\build.ps1

Write-Host "🔨 Compilando Calculadora Distribuída (gRPC)" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green

# Verificar se Go está instalado
try {
    $goVersion = go version
    Write-Host "✅ Go encontrado: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go não encontrado. Instale Go 1.21+ primeiro." -ForegroundColor Red
    exit 1
}

# Verificar se protoc está instalado
try {
    $protocVersion = protoc --version
    Write-Host "✅ protoc encontrado: $protocVersion" -ForegroundColor Green
} catch {
    Write-Host "⚠️  protoc não encontrado. Execute scripts\install_protoc.ps1 ou consulte SETUP.md" -ForegroundColor Yellow
    $response = Read-Host "Continuar sem gerar código a partir do .proto? (y/n)"
    if ($response -ne "y") {
        exit 1
    }
}

# Criar diretório bin
Write-Host "`n📁 Criando diretório bin..." -ForegroundColor Cyan
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Baixar dependências
Write-Host "`n📦 Baixando dependências..." -ForegroundColor Cyan
go mod download
go mod tidy

# Gerar código a partir do .proto (se protoc disponível)
if (Get-Command "protoc" -ErrorAction SilentlyContinue) {
    Write-Host "`n🔧 Gerando código a partir do .proto..." -ForegroundColor Cyan
    protoc --go_out=. --go_opt=paths=source_relative `
        --go-grpc_out=. --go-grpc_opt=paths=source_relative `
        proto/calculator.proto

    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Código gerado com sucesso!" -ForegroundColor Green
    } else {
        Write-Host "❌ Erro ao gerar código" -ForegroundColor Red
        exit 1
    }
}

# Compilar binários
Write-Host "`n🔨 Compilando binários..." -ForegroundColor Cyan

$components = @(
    @{Name="Add Server"; Path="./cmd/grpc_add_server"; Output="bin/grpc_add_server.exe"},
    @{Name="Subtract Server"; Path="./cmd/grpc_sub_server"; Output="bin/grpc_sub_server.exe"},
    @{Name="Multiply Server"; Path="./cmd/grpc_mult_server"; Output="bin/grpc_mult_server.exe"},
    @{Name="Divide Server"; Path="./cmd/grpc_div_server"; Output="bin/grpc_div_server.exe"},
    @{Name="Dispatcher"; Path="./cmd/grpc_dispatcher"; Output="bin/grpc_dispatcher.exe"},
    @{Name="Client"; Path="./cmd/grpc_client"; Output="bin/grpc_client.exe"}
)

$success = $true
foreach ($comp in $components) {
    Write-Host "  ➤ Compilando $($comp.Name)..." -ForegroundColor White
    go build -o $comp.Output $comp.Path

    if ($LASTEXITCODE -eq 0) {
        Write-Host "    ✅ $($comp.Name) compilado" -ForegroundColor Green
    } else {
        Write-Host "    ❌ Erro ao compilar $($comp.Name)" -ForegroundColor Red
        $success = $false
    }
}

Write-Host "`n=============================================" -ForegroundColor Green
if ($success) {
    Write-Host "✅ Build concluído com sucesso!" -ForegroundColor Green
    Write-Host "`nExecute o sistema com:" -ForegroundColor Cyan
    Write-Host "  powershell -ExecutionPolicy Bypass -File scripts\run.ps1" -ForegroundColor White
} else {
    Write-Host "❌ Build falhou. Verifique os erros acima." -ForegroundColor Red
    exit 1
}

# Script de build para Windows
# Execute: powershell -ExecutionPolicy Bypass -File scripts\build.ps1

Write-Host "Compilando Calculadora Distribuida (gRPC)" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green

# Verificar se Go esta instalado
try {
    $goVersion = go version
    Write-Host "Go encontrado: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "Go nao encontrado. Instale Go 1.21+ primeiro." -ForegroundColor Red
    exit 1
}

# Criar diretorio bin
Write-Host "`nCriando diretorio bin..." -ForegroundColor Cyan
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Baixar dependencias
Write-Host "`nBaixando dependencias..." -ForegroundColor Cyan
go mod download
go mod tidy

# Gerar codigo a partir do .proto
Write-Host "`nGerando codigo a partir do .proto..." -ForegroundColor Cyan
if (Get-Command "protoc" -ErrorAction SilentlyContinue) {
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/calculator.proto

    if ($LASTEXITCODE -eq 0) {
        Write-Host "Codigo gerado com sucesso!" -ForegroundColor Green
    } else {
        Write-Host "Erro ao gerar codigo" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "protoc nao encontrado, pulando geracao de codigo proto" -ForegroundColor Yellow
}

# Compilar binarios
Write-Host "`nCompilando binarios..." -ForegroundColor Cyan

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
    Write-Host "  Compilando $($comp.Name)..." -ForegroundColor White
    go build -o $comp.Output $comp.Path

    if ($LASTEXITCODE -eq 0) {
        Write-Host "    $($comp.Name) compilado" -ForegroundColor Green
    } else {
        Write-Host "    Erro ao compilar $($comp.Name)" -ForegroundColor Red
        $success = $false
    }
}

Write-Host "`n=============================================" -ForegroundColor Green
if ($success) {
    Write-Host "Build concluido com sucesso!" -ForegroundColor Green
    Write-Host "`nExecute o sistema com:" -ForegroundColor Cyan
    Write-Host "  .\scripts\run.ps1" -ForegroundColor White
} else {
    Write-Host "Build falhou. Verifique os erros acima." -ForegroundColor Red
    exit 1
}

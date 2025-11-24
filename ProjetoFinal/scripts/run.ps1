# Script para executar todos os componentes do sistema
# Execute: powershell -ExecutionPolicy Bypass -File scripts\run.ps1

Write-Host "Iniciando Sistema de Calculadora Distribuida (gRPC)" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

# Verificar se os binarios existem
$binaries = @(
    "bin\grpc_add_server.exe",
    "bin\grpc_sub_server.exe",
    "bin\grpc_mult_server.exe",
    "bin\grpc_div_server.exe",
    "bin\grpc_dispatcher.exe",
    "bin\grpc_client.exe"
)

$allExist = $true
foreach ($bin in $binaries) {
    if (-not (Test-Path $bin)) {
        Write-Host "Binario nao encontrado: $bin" -ForegroundColor Red
        $allExist = $false
    }
}

if (-not $allExist) {
    Write-Host "`nCompile o projeto primeiro com: go build" -ForegroundColor Yellow
    exit 1
}

# Funcao para matar processos anteriores
function Stop-OldProcesses {
    Write-Host "`nLimpando processos anteriores..." -ForegroundColor Yellow
    $processes = @(
        "grpc_add_server",
        "grpc_sub_server",
        "grpc_mult_server",
        "grpc_div_server",
        "grpc_dispatcher"
    )

    foreach ($proc in $processes) {
        Get-Process -Name $proc -ErrorAction SilentlyContinue | Stop-Process -Force
    }
    Start-Sleep -Seconds 1
}

# Limpar processos anteriores
Stop-OldProcesses

# Iniciar servidores de operacao
Write-Host "`nIniciando servidores de operacao..." -ForegroundColor Cyan

Write-Host "  Add Server (porta 50052)" -ForegroundColor White
Start-Process -NoNewWindow -FilePath .\bin\grpc_add_server.exe

Write-Host "  Subtract Server (porta 50053)" -ForegroundColor White
Start-Process -NoNewWindow -FilePath .\bin\grpc_sub_server.exe

Write-Host "  Multiply Server (porta 50054)" -ForegroundColor White
Start-Process -NoNewWindow -FilePath .\bin\grpc_mult_server.exe

Write-Host "  Divide Server (porta 50055)" -ForegroundColor White
Start-Process -NoNewWindow -FilePath .\bin\grpc_div_server.exe

# Aguardar servidores iniciarem
Write-Host "`nAguardando servidores iniciarem..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Iniciar dispatcher
Write-Host "`nIniciando Dispatcher (porta 50051)..." -ForegroundColor Cyan
Start-Process -NoNewWindow -FilePath .\bin\grpc_dispatcher.exe

# Aguardar dispatcher iniciar
Write-Host "`nAguardando dispatcher iniciar..." -ForegroundColor Yellow
Start-Sleep -Seconds 2

# Iniciar cliente
Write-Host "`nIniciando Cliente..." -ForegroundColor Cyan
Write-Host "=================================================" -ForegroundColor Green
Write-Host ""

& .\bin\grpc_client.exe

# Quando o cliente encerrar, parar todos os processos
Write-Host "`nEncerrando todos os servidores..." -ForegroundColor Yellow
Stop-OldProcesses

Write-Host "`nSistema encerrado!" -ForegroundColor Green

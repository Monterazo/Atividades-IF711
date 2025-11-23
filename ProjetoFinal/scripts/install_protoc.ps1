# Script de instalação do protoc para Windows
# Execute como administrador: powershell -ExecutionPolicy Bypass -File install_protoc.ps1

param(
    [string]$Version = "25.1"
)

Write-Host "Instalando Protocol Buffers Compiler (protoc) v$Version" -ForegroundColor Green

# Configurações
$ProtocVersion = $Version
$ProtocUrl = "https://github.com/protocolbuffers/protobuf/releases/download/v$ProtocVersion/protoc-$ProtocVersion-win64.zip"
$InstallDir = "C:\protoc"
$ZipFile = "$env:TEMP\protoc.zip"

# Criar diretório de instalação
Write-Host "Criando diretório de instalação em $InstallDir..." -ForegroundColor Yellow
if (Test-Path $InstallDir) {
    Remove-Item -Recurse -Force $InstallDir
}
New-Item -ItemType Directory -Path $InstallDir | Out-Null

# Baixar protoc
Write-Host "Baixando protoc v$ProtocVersion..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri $ProtocUrl -OutFile $ZipFile
} catch {
    Write-Host "Erro ao baixar protoc: $_" -ForegroundColor Red
    exit 1
}

# Extrair
Write-Host "Extraindo arquivos..." -ForegroundColor Yellow
Expand-Archive -Path $ZipFile -DestinationPath $InstallDir -Force

# Adicionar ao PATH
Write-Host "Adicionando ao PATH do sistema..." -ForegroundColor Yellow
$ProtocBinPath = "$InstallDir\bin"
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")

if ($CurrentPath -notlike "*$ProtocBinPath*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$CurrentPath;$ProtocBinPath",
        "Machine"
    )
    Write-Host "✅ Adicionado ao PATH!" -ForegroundColor Green
} else {
    Write-Host "✅ Já está no PATH!" -ForegroundColor Green
}

# Limpar
Remove-Item $ZipFile -Force

# Verificar instalação
Write-Host "`nVerificando instalação..." -ForegroundColor Yellow
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine")
& "$ProtocBinPath\protoc.exe" --version

Write-Host "`n✅ Instalação concluída!" -ForegroundColor Green
Write-Host "⚠️  IMPORTANTE: Reinicie o terminal para que as mudanças no PATH tenham efeito." -ForegroundColor Yellow
Write-Host "`nPróximos passos:" -ForegroundColor Cyan
Write-Host "1. Feche e reabra o terminal" -ForegroundColor White
Write-Host "2. Execute: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" -ForegroundColor White
Write-Host "3. Execute: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" -ForegroundColor White

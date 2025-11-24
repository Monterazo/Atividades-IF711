@echo off
REM Script de Benchmark para gRPC
REM Uso: benchmark-grpc.bat [expressao] [clientes] [requisicoes]

setlocal

set EXPRESSION=%1
set CLIENTS=%2
set REQUESTS=%3

if "%EXPRESSION%"=="" set EXPRESSION=((4+3)*2)/5
if "%CLIENTS%"=="" set CLIENTS=10
if "%REQUESTS%"=="" set REQUESTS=100

echo.
echo ============================================
echo   BENCHMARK gRPC - Calculadora Distribuida
echo ============================================
echo   Expressao: %EXPRESSION%
echo   Clientes:  %CLIENTS%
echo   Requests:  %REQUESTS%
echo ============================================
echo.

echo Compilando benchmark...
go build -o bin\grpc_benchmark.exe cmd\benchmark\grpc_benchmark.go

if %ERRORLEVEL% NEQ 0 (
    echo Erro ao compilar benchmark!
    exit /b 1
)

echo Executando benchmark...
echo.

bin\grpc_benchmark.exe -expr="%EXPRESSION%" -clients=%CLIENTS% -reqs=%REQUESTS%

endlocal

.PHONY: proto build run-servers run-dispatcher run-client clean

# Gera código a partir do arquivo .proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/calculator.proto

# Baixa dependências
deps:
	go mod download
	go mod tidy

# Compila todos os binários
build: proto
	@echo "Compilando servidores..."
	go build -o bin/grpc_add_server ./cmd/grpc_add_server
	go build -o bin/grpc_sub_server ./cmd/grpc_sub_server
	go build -o bin/grpc_mult_server ./cmd/grpc_mult_server
	go build -o bin/grpc_div_server ./cmd/grpc_div_server
	go build -o bin/grpc_dispatcher ./cmd/grpc_dispatcher
	go build -o bin/grpc_client ./cmd/grpc_client
	@echo "Compilação concluída!"

# Executa todos os servidores em background
run-servers:
	@echo "Iniciando servidores de operação..."
	start /B bin/grpc_add_server.exe
	start /B bin/grpc_sub_server.exe
	start /B bin/grpc_mult_server.exe
	start /B bin/grpc_div_server.exe
	@echo "Servidores iniciados!"

# Executa o dispatcher
run-dispatcher:
	@echo "Iniciando dispatcher..."
	bin/grpc_dispatcher.exe

# Executa o cliente
run-client:
	@echo "Iniciando cliente..."
	bin/grpc_client.exe

# Limpa binários
clean:
	rm -rf bin/
	@echo "Binários removidos!"

# Executa tudo (servidores + dispatcher + client)
run-all: build
	@echo "Iniciando todos os componentes..."
	@echo "Pressione Ctrl+C para encerrar"
	start /B bin/grpc_add_server.exe
	start /B bin/grpc_sub_server.exe
	start /B bin/grpc_mult_server.exe
	start /B bin/grpc_div_server.exe
	timeout /t 2 /nobreak > nul
	start /B bin/grpc_dispatcher.exe
	timeout /t 2 /nobreak > nul
	bin/grpc_client.exe

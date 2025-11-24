.PHONY: proto build build-grpc build-rabbitmq build-benchmark run-servers run-dispatcher run-client clean benchmark-grpc benchmark-rabbitmq

# Gera código a partir do arquivo .proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/calculator.proto

# Baixa dependências
deps:
	go mod download
	go mod tidy

# Compila todos os binários (gRPC e RabbitMQ)
build: build-grpc build-rabbitmq

# Compila binários gRPC
build-grpc: proto
	@echo "Compilando servidores gRPC..."
	go build -o bin/grpc_add_server ./cmd/grpc_add_server
	go build -o bin/grpc_sub_server ./cmd/grpc_sub_server
	go build -o bin/grpc_mult_server ./cmd/grpc_mult_server
	go build -o bin/grpc_div_server ./cmd/grpc_div_server
	go build -o bin/grpc_dispatcher ./cmd/grpc_dispatcher
	go build -o bin/grpc_client ./cmd/grpc_client
	@echo "Compilação gRPC concluída!"

# Compila binários RabbitMQ
build-rabbitmq:
	@echo "Compilando servidores RabbitMQ..."
	go build -o bin/rabbitmq_add_server ./cmd/rabbitmq_add_server
	go build -o bin/rabbitmq_sub_server ./cmd/rabbitmq_sub_server
	go build -o bin/rabbitmq_mult_server ./cmd/rabbitmq_mult_server
	go build -o bin/rabbitmq_div_server ./cmd/rabbitmq_div_server
	go build -o bin/rabbitmq_dispatcher ./cmd/rabbitmq_dispatcher
	go build -o bin/rabbitmq_client ./cmd/rabbitmq_client
	@echo "Compilação RabbitMQ concluída!"

# Compila ferramentas de benchmark
build-benchmark:
	@echo "Compilando ferramentas de benchmark..."
	go build -o bin/grpc_benchmark.exe ./cmd/benchmark/grpc_benchmark.go
	go build -o bin/rabbitmq_benchmark.exe ./cmd/benchmark/rabbitmq_benchmark.go
	@echo "Compilação de benchmarks concluída!"

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

# Executa tudo gRPC (servidores + dispatcher + client)
run-all: build-grpc
	@echo "Iniciando todos os componentes gRPC..."
	@echo "Pressione Ctrl+C para encerrar"
	start /B bin/grpc_add_server.exe
	start /B bin/grpc_sub_server.exe
	start /B bin/grpc_mult_server.exe
	start /B bin/grpc_div_server.exe
	timeout /t 2 /nobreak > nul
	start /B bin/grpc_dispatcher.exe
	timeout /t 2 /nobreak > nul
	bin/grpc_client.exe

# Executa tudo RabbitMQ (servidores + dispatcher + client)
run-all-rabbitmq: build-rabbitmq
	@echo "Iniciando todos os componentes RabbitMQ..."
	@echo "IMPORTANTE: Certifique-se de que o RabbitMQ está rodando!"
	@echo "Pressione Ctrl+C para encerrar"
	start /B bin/rabbitmq_add_server.exe
	start /B bin/rabbitmq_sub_server.exe
	start /B bin/rabbitmq_mult_server.exe
	start /B bin/rabbitmq_div_server.exe
	timeout /t 2 /nobreak > nul
	start /B bin/rabbitmq_dispatcher.exe
	timeout /t 2 /nobreak > nul
	bin/rabbitmq_client.exe

# Executa benchmark gRPC
# Uso: make benchmark-grpc EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100
benchmark-grpc: build-benchmark
	@echo "Executando benchmark gRPC..."
	./bin/grpc_benchmark.exe -expr="$(or $(EXPR),((4+3)*2)/5)" -clients=$(or $(CLIENTS),10) -reqs=$(or $(REQS),100)

# Executa benchmark RabbitMQ
# Uso: make benchmark-rabbitmq EXPR="((4+3)*2)/5" CLIENTS=10 REQS=100
benchmark-rabbitmq: build-benchmark
	@echo "Executando benchmark RabbitMQ..."
	@echo "IMPORTANTE: Certifique-se de que o RabbitMQ e os servidores estão rodando!"
	./bin/rabbitmq_benchmark.exe -expr="$(or $(EXPR),((4+3)*2)/5)" -clients=$(or $(CLIENTS),10) -reqs=$(or $(REQS),100)

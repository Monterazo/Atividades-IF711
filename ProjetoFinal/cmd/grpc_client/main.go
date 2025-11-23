package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/Monterazo/Atividades-IF711/ProjetoFinal/proto"
	"google.golang.org/grpc"
)

const (
	dispatcherAddr = "localhost:50051"
	defaultTimeout = 30000 // 30 segundos em milissegundos
)

func main() {
	log.Println("Cliente Calculadora gRPC")
	log.Println("========================")

	// Conecta ao dispatcher
	conn, err := grpc.Dial(dispatcherAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Falha ao conectar ao dispatcher: %v", err)
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)
	log.Printf("Conectado ao dispatcher em %s\n", dispatcherAddr)

	// Loop de interação
	scanner := bufio.NewScanner(os.Stdin)
	expressionCounter := 0

	fmt.Println("\nDigite uma expressão matemática (ou 'sair' para encerrar):")
	fmt.Println("Exemplos: ((4+3)*2)/5, 10+20*3, (15-5)/2")

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "sair" || strings.ToLower(input) == "exit" {
			fmt.Println("Encerrando cliente...")
			break
		}

		// Prepara requisição
		expressionCounter++
		expressionID := fmt.Sprintf("expr_%d", expressionCounter)

		req := &pb.ExpressionRequest{
			ExpressionId: expressionID,
			Expression:   input,
			DeadlineMs:   defaultTimeout,
		}

		// Cria contexto com timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Millisecond)

		log.Printf("Enviando expressão: %s", input)

		// Envia requisição
		startTime := time.Now()
		resp, err := client.Calculate(ctx, req)
		duration := time.Since(startTime)
		cancel()

		if err != nil {
			fmt.Printf("❌ Erro ao calcular: %v\n", err)
			log.Printf("Erro: %v", err)
			continue
		}

		// Processa resposta
		if resp.Error != nil {
			fmt.Printf("❌ Erro: [%s] %s\n", resp.Error.Code, resp.Error.Message)
			log.Printf("Erro retornado: [%s] %s", resp.Error.Code, resp.Error.Message)
		} else {
			fmt.Printf("✅ Resultado: %s = %f\n", input, resp.Result)
			fmt.Printf("⏱️  Tempo de execução: %v\n", duration)
			log.Printf("Resultado: %f (tempo: %v)", resp.Result, duration)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erro ao ler entrada: %v", err)
	}
}

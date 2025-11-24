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
	// Gera ID único para este cliente
	clientID := fmt.Sprintf("CLIENT-%d", time.Now().Unix()%10000)

	log.Printf("[%s] Cliente Calculadora gRPC", clientID)
	log.Printf("[%s] ========================", clientID)

	// Conecta ao dispatcher
	conn, err := grpc.Dial(dispatcherAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("[%s] Falha ao conectar ao dispatcher: %v", clientID, err)
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)
	log.Printf("[%s] Conectado ao dispatcher em %s\n", clientID, dispatcherAddr)

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
		expressionID := fmt.Sprintf("%s_expr_%d", clientID, expressionCounter)

		req := &pb.ExpressionRequest{
			ExpressionId: expressionID,
			Expression:   input,
			DeadlineMs:   defaultTimeout,
		}

		// Cria contexto com timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Millisecond)

		log.Printf("[%s] Enviando expressão: %s (ID: %s)", clientID, input, expressionID)

		// Envia requisição
		startTime := time.Now()
		resp, err := client.Calculate(ctx, req)
		duration := time.Since(startTime)
		cancel()

		if err != nil {
			fmt.Printf("❌ Erro ao calcular: %v\n", err)
			log.Printf("[%s] Erro: %v", clientID, err)
			continue
		}

		// Processa resposta
		if resp.Error != nil {
			fmt.Printf("❌ Erro: [%s] %s\n", resp.Error.Code, resp.Error.Message)
			log.Printf("[%s] Erro retornado: [%s] %s", clientID, resp.Error.Code, resp.Error.Message)
		} else {
			fmt.Printf("✅ Resultado: %s = %f\n", input, resp.Result)
			fmt.Printf("⏱️  Tempo de execução: %v\n", duration)
			log.Printf("[%s] Resultado: %f (tempo: %v)", clientID, resp.Result, duration)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erro ao ler entrada: %v", err)
	}
}

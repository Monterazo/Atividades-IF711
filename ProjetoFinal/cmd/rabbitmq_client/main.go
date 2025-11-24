package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/rabbitmq"
)

const (
	rabbitmqURL    = "amqp://guest:guest@localhost:5672/"
	defaultTimeout = 30000 // 30 segundos em milissegundos
)

func main() {
	// Gera ID único para este cliente
	clientID := fmt.Sprintf("CLIENT-%d", time.Now().Unix()%10000)

	log.Printf("[%s] Cliente Calculadora RabbitMQ", clientID)
	log.Printf("[%s] ========================", clientID)

	// Conecta ao RabbitMQ
	conn, err := rabbitmq.NewConnection(rabbitmqURL)
	if err != nil {
		log.Fatalf("[%s] Erro ao conectar ao RabbitMQ: %v", clientID, err)
	}
	defer conn.Close()

	// Declara filas necessárias
	if err := conn.DeclareQueue(rabbitmq.RequestQueue); err != nil {
		log.Fatalf("[%s] Erro ao declarar fila de requests: %v", clientID, err)
	}
	if err := conn.DeclareQueue(rabbitmq.ResponseQueue); err != nil {
		log.Fatalf("[%s] Erro ao declarar fila de responses: %v", clientID, err)
	}

	log.Printf("[%s] Conectado ao RabbitMQ em %s", clientID, rabbitmqURL)

	// Inicia consumidor de respostas
	msgs, err := conn.Consume(rabbitmq.ResponseQueue)
	if err != nil {
		log.Fatalf("[%s] Erro ao consumir fila de responses: %v", clientID, err)
	}

	// Canal para respostas
	responseChan := make(chan rabbitmq.ExpressionResponse, 10)

	// Goroutine para processar respostas
	go func() {
		for msg := range msgs {
			var resp rabbitmq.ExpressionResponse
			if err := json.Unmarshal(msg.Body, &resp); err != nil {
				log.Printf("[%s] Erro ao decodificar resposta: %v", clientID, err)
				msg.Nack(false, false)
				continue
			}

			// Filtra apenas respostas para este cliente
			if strings.HasPrefix(resp.ExpressionID, clientID) {
				responseChan <- resp
				msg.Ack(false)
			} else {
				// Rejeita mensagem que não é para este cliente
				msg.Nack(false, true)
			}
		}
	}()

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

		req := rabbitmq.ExpressionRequest{
			ExpressionID: expressionID,
			Expression:   input,
			DeadlineMs:   defaultTimeout,
		}

		// Serializa requisição
		reqBytes, err := json.Marshal(req)
		if err != nil {
			fmt.Printf("❌ Erro ao serializar requisição: %v\n", err)
			continue
		}

		log.Printf("[%s] Enviando expressão: %s (ID: %s)", clientID, input, expressionID)

		// Envia requisição
		startTime := time.Now()
		if err := conn.Publish(rabbitmq.RequestQueue, reqBytes); err != nil {
			fmt.Printf("❌ Erro ao enviar requisição: %v\n", err)
			log.Printf("[%s] Erro ao enviar: %v", clientID, err)
			continue
		}

		// Aguarda resposta com timeout
		select {
		case resp := <-responseChan:
			duration := time.Since(startTime)

			if resp.Error != nil {
				fmt.Printf("❌ Erro: [%s] %s\n", resp.Error.Code, resp.Error.Message)
				log.Printf("[%s] Erro retornado: [%s] %s", clientID, resp.Error.Code, resp.Error.Message)
			} else {
				fmt.Printf("✅ Resultado: %s = %f\n", input, resp.Result)
				fmt.Printf("⏱️  Tempo de execução: %v\n", duration)
				log.Printf("[%s] Resultado: %f (tempo: %v)", clientID, resp.Result, duration)
			}

		case <-time.After(time.Duration(defaultTimeout) * time.Millisecond):
			fmt.Println("❌ Timeout ao aguardar resposta")
			log.Printf("[%s] Timeout ao aguardar resposta para %s", clientID, expressionID)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Erro ao ler entrada: %v", err)
	}
}

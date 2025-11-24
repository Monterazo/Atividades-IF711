package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/rabbitmq"
)

const (
	rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	operation   = "multiply"
)

func main() {
	serverName := strings.ToUpper(operation)
	log.Printf("[%s] Iniciando servidor de operação: %s", serverName, operation)

	// Conecta ao RabbitMQ
	conn, err := rabbitmq.NewConnection(rabbitmqURL)
	if err != nil {
		log.Fatalf("[%s] Erro ao conectar ao RabbitMQ: %v", serverName, err)
	}
	defer conn.Close()

	// Declara fila de operação
	queue := rabbitmq.GetQueueForOperation(operation)
	if err := conn.DeclareQueue(queue); err != nil {
		log.Fatalf("[%s] Erro ao declarar fila: %v", serverName, err)
	}

	// Declara fila de resultados
	if err := conn.DeclareQueue(rabbitmq.ResultsQueue); err != nil {
		log.Fatalf("[%s] Erro ao declarar fila de resultados: %v", serverName, err)
	}

	// Consome mensagens da fila
	msgs, err := conn.Consume(queue)
	if err != nil {
		log.Fatalf("[%s] Erro ao consumir fila: %v", serverName, err)
	}

	log.Printf("[%s] Servidor %s pronto para receber operações", serverName, operation)

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var req rabbitmq.OperationRequest
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				log.Printf("[%s] Erro ao decodificar requisição: %v", serverName, err)
				msg.Nack(false, false)
				continue
			}

			// Extrai clientID
			clientID := "UNKNOWN"
			parts := strings.Split(req.ExpressionID, "_expr_")
			if len(parts) > 0 {
				clientID = parts[0]
			}

			log.Printf("[%s] [%s] Recebida operação: %s(%v) [Step: %s]", serverName, clientID, req.Operation, req.Numbers, req.StepID)

			// Valida operação
			if req.Operation != operation {
				log.Printf("[%s] [%s] Operação inválida: esperado %s, recebido %s", serverName, clientID, operation, req.Operation)
				resp := rabbitmq.OperationResponse{
					ExpressionID: req.ExpressionID,
					StepID:       req.StepID,
					Error: &rabbitmq.ErrorInfo{
						Code:    "INVALID_OPERATION",
						Message: "Este servidor só processa operações " + operation,
					},
				}
				sendResponse(conn, resp, serverName)
				msg.Ack(false)
				continue
			}

			// Executa operação
			result, err := rabbitmq.ExecuteOperation(req.Operation, req.Numbers)
			if err != nil {
				log.Printf("[%s] [%s] Erro ao executar operação: %v", serverName, clientID, err)
				errorCode := "EXECUTION_ERROR"
				if err.Error() == "divisão por zero" {
					errorCode = "DIV_BY_ZERO"
				}
				resp := rabbitmq.OperationResponse{
					ExpressionID: req.ExpressionID,
					StepID:       req.StepID,
					Error: &rabbitmq.ErrorInfo{
						Code:    errorCode,
						Message: err.Error(),
					},
				}
				sendResponse(conn, resp, serverName)
				msg.Ack(false)
				continue
			}

			log.Printf("[%s] [%s] Operação executada com sucesso: %f", serverName, clientID, result)

			// Envia resultado
			resp := rabbitmq.OperationResponse{
				ExpressionID: req.ExpressionID,
				StepID:       req.StepID,
				Result:       result,
			}
			sendResponse(conn, resp, serverName)
			msg.Ack(false)
		}
	}()

	<-forever
}

func sendResponse(conn *rabbitmq.Connection, resp rabbitmq.OperationResponse, serverName string) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[%s] Erro ao serializar resposta: %v", serverName, err)
		return
	}

	if err := conn.Publish(rabbitmq.ResultsQueue, respBytes); err != nil {
		log.Printf("[%s] Erro ao publicar resposta: %v", serverName, err)
	}
}

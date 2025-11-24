package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/core"
	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/rabbitmq"
)

const (
	rabbitmqURL = "amqp://guest:guest@localhost:5672/"
)

type PendingStep struct {
	ExpressionID string
	TotalSteps   int
	Results      map[string]float64
	Steps        []core.Step
	Mutex        sync.Mutex
	ResponseSent bool
}

type Dispatcher struct {
	conn          *rabbitmq.Connection
	parser        *core.Parser
	pendingSteps  map[string]*PendingStep
	pendingMutex  sync.RWMutex
}

func NewDispatcher(conn *rabbitmq.Connection) *Dispatcher {
	return &Dispatcher{
		conn:         conn,
		parser:       core.NewParser(),
		pendingSteps: make(map[string]*PendingStep),
	}
}

func (d *Dispatcher) processRequest(msg []byte) {
	var req rabbitmq.ExpressionRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Printf("[DISPATCHER] Erro ao decodificar requisição: %v", err)
		return
	}

	// Extrai clientID
	clientID := "UNKNOWN"
	parts := strings.Split(req.ExpressionID, "_expr_")
	if len(parts) > 0 {
		clientID = parts[0]
	}

	log.Printf("[DISPATCHER] [%s] Recebida expressão: %s (ID: %s)", clientID, req.Expression, req.ExpressionID)

	// Parse da expressão
	steps, rpnStr, err := d.parser.ParseWithRPN(req.Expression)
	if err != nil {
		log.Printf("[DISPATCHER] [%s] Erro ao fazer parse: %v", clientID, err)
		d.sendErrorResponse(req.ExpressionID, "PARSE_ERROR", fmt.Sprintf("Erro ao fazer parse: %v", err))
		return
	}

	log.Printf("[DISPATCHER] [%s] RPN: %s", clientID, rpnStr)
	log.Printf("[DISPATCHER] [%s] Expressão parseada em %d steps", clientID, len(steps))

	// Registra expressão pendente
	d.pendingMutex.Lock()
	d.pendingSteps[req.ExpressionID] = &PendingStep{
		ExpressionID: req.ExpressionID,
		TotalSteps:   len(steps),
		Results:      make(map[string]float64),
		Steps:        steps,
		ResponseSent: false,
	}
	d.pendingMutex.Unlock()

	// Processa primeiro step
	d.processNextStep(req.ExpressionID, clientID, req.DeadlineMs)
}

func (d *Dispatcher) processNextStep(expressionID, clientID string, deadlineMs int64) {
	d.pendingMutex.RLock()
	pending, exists := d.pendingSteps[expressionID]
	d.pendingMutex.RUnlock()

	if !exists {
		return
	}

	pending.Mutex.Lock()
	defer pending.Mutex.Unlock()

	currentStepIndex := len(pending.Results)
	if currentStepIndex >= pending.TotalSteps {
		return
	}

	step := pending.Steps[currentStepIndex]

	// Substitui referências a resultados anteriores
	numbers := make([]float64, len(step.Numbers))
	copy(numbers, step.Numbers)

	for _, dep := range step.DependsOn {
		// Extrai a parte após "result_" (ex: "result_step0" -> "step0")
		refSuffix := dep.Reference[len("result_"):]
		// Constrói o stepID completo (ex: "CLIENT-9219_expr_2_step0")
		stepID := fmt.Sprintf("%s_%s", expressionID, refSuffix)
		if result, ok := pending.Results[stepID]; ok {
			numbers[dep.Position] = result
		}
	}

	// Prepara requisição de operação
	stepID := fmt.Sprintf("%s_step%d", expressionID, currentStepIndex)
	opReq := rabbitmq.OperationRequest{
		ExpressionID: expressionID,
		StepID:       stepID,
		Operation:    step.Operation,
		Numbers:      numbers,
		DeadlineMs:   deadlineMs,
	}

	log.Printf("[DISPATCHER] [%s] Executando step %d: %s(%v)", clientID, currentStepIndex, step.Operation, numbers)

	// Serializa e envia para a fila da operação
	opReqBytes, err := json.Marshal(opReq)
	if err != nil {
		log.Printf("[DISPATCHER] [%s] Erro ao serializar operação: %v", clientID, err)
		d.sendErrorResponse(expressionID, "SERIALIZATION_ERROR", fmt.Sprintf("Erro ao serializar: %v", err))
		return
	}

	queue := rabbitmq.GetQueueForOperation(step.Operation)
	if queue == "" {
		log.Printf("[DISPATCHER] [%s] Operação desconhecida: %s", clientID, step.Operation)
		d.sendErrorResponse(expressionID, "UNKNOWN_OPERATION", fmt.Sprintf("Operação desconhecida: %s", step.Operation))
		return
	}

	if err := d.conn.Publish(queue, opReqBytes); err != nil {
		log.Printf("[DISPATCHER] [%s] Erro ao publicar operação: %v", clientID, err)
		d.sendErrorResponse(expressionID, "PUBLISH_ERROR", fmt.Sprintf("Erro ao publicar: %v", err))
		return
	}
}

func (d *Dispatcher) processOperationResult(msg []byte) {
	var resp rabbitmq.OperationResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		log.Printf("[DISPATCHER] Erro ao decodificar resultado: %v", err)
		return
	}

	// Extrai clientID
	clientID := "UNKNOWN"
	parts := strings.Split(resp.ExpressionID, "_expr_")
	if len(parts) > 0 {
		clientID = parts[0]
	}

	log.Printf("[DISPATCHER] [%s] Recebido resultado do step %s", clientID, resp.StepID)

	// Verifica se houve erro
	if resp.Error != nil {
		log.Printf("[DISPATCHER] [%s] Erro no step: %s - %s", clientID, resp.Error.Code, resp.Error.Message)
		d.sendErrorResponse(resp.ExpressionID, resp.Error.Code, resp.Error.Message)
		d.cleanupExpression(resp.ExpressionID)
		return
	}

	d.pendingMutex.RLock()
	pending, exists := d.pendingSteps[resp.ExpressionID]
	d.pendingMutex.RUnlock()

	if !exists {
		log.Printf("[DISPATCHER] [%s] Expressão não encontrada: %s", clientID, resp.ExpressionID)
		return
	}

	pending.Mutex.Lock()
	pending.Results[resp.StepID] = resp.Result
	currentStepCount := len(pending.Results)
	pending.Mutex.Unlock()

	log.Printf("[DISPATCHER] [%s] Step completado: resultado = %f (%d/%d)", clientID, resp.Result, currentStepCount, pending.TotalSteps)

	// Verifica se todos os steps foram completados
	if currentStepCount >= pending.TotalSteps {
		// Expressão completa
		log.Printf("[DISPATCHER] [%s] Expressão calculada com sucesso: %f", clientID, resp.Result)
		d.sendSuccessResponse(resp.ExpressionID, resp.Result)
		d.cleanupExpression(resp.ExpressionID)
	} else {
		// Processa próximo step
		// Pequeno delay para garantir processamento sequencial
		time.Sleep(10 * time.Millisecond)
		d.processNextStep(resp.ExpressionID, clientID, 5000)
	}
}

func (d *Dispatcher) sendSuccessResponse(expressionID string, result float64) {
	resp := rabbitmq.ExpressionResponse{
		ExpressionID: expressionID,
		Result:       result,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[DISPATCHER] Erro ao serializar resposta: %v", err)
		return
	}

	if err := d.conn.Publish(rabbitmq.ResponseQueue, respBytes); err != nil {
		log.Printf("[DISPATCHER] Erro ao publicar resposta: %v", err)
	}
}

func (d *Dispatcher) sendErrorResponse(expressionID, code, message string) {
	d.pendingMutex.RLock()
	pending, exists := d.pendingSteps[expressionID]
	d.pendingMutex.RUnlock()

	if exists {
		pending.Mutex.Lock()
		if pending.ResponseSent {
			pending.Mutex.Unlock()
			return
		}
		pending.ResponseSent = true
		pending.Mutex.Unlock()
	}

	resp := rabbitmq.ExpressionResponse{
		ExpressionID: expressionID,
		Error: &rabbitmq.ErrorInfo{
			Code:    code,
			Message: message,
		},
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[DISPATCHER] Erro ao serializar resposta de erro: %v", err)
		return
	}

	if err := d.conn.Publish(rabbitmq.ResponseQueue, respBytes); err != nil {
		log.Printf("[DISPATCHER] Erro ao publicar resposta de erro: %v", err)
	}
}

func (d *Dispatcher) cleanupExpression(expressionID string) {
	d.pendingMutex.Lock()
	delete(d.pendingSteps, expressionID)
	d.pendingMutex.Unlock()
}

func main() {
	log.Println("Iniciando Dispatcher RabbitMQ...")

	// Conecta ao RabbitMQ
	conn, err := rabbitmq.NewConnection(rabbitmqURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Configura filas
	if err := rabbitmq.SetupQueues(conn); err != nil {
		log.Fatalf("Erro ao configurar filas: %v", err)
	}

	log.Println("Filas configuradas com sucesso")

	// Cria dispatcher
	dispatcher := NewDispatcher(conn)

	// Consome requisições
	requests, err := conn.Consume(rabbitmq.RequestQueue)
	if err != nil {
		log.Fatalf("Erro ao consumir fila de requests: %v", err)
	}

	// Consome resultados de operações
	results, err := conn.Consume(rabbitmq.ResultsQueue)
	if err != nil {
		log.Fatalf("Erro ao consumir fila de results: %v", err)
	}

	log.Println("Dispatcher pronto para receber requisições")

	// Processa mensagens
	forever := make(chan bool)

	go func() {
		for msg := range requests {
			dispatcher.processRequest(msg.Body)
			msg.Ack(false)
		}
	}()

	go func() {
		for msg := range results {
			dispatcher.processOperationResult(msg.Body)
			msg.Ack(false)
		}
	}()

	<-forever
}

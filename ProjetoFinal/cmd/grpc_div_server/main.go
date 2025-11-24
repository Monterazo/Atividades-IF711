package main

import (
	"context"
	"log"
	"net"
	"strings"

	pb "github.com/Monterazo/Atividades-IF711/ProjetoFinal/proto"
	grpcOps "github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/grpc"
	"google.golang.org/grpc"
)

const (
	port = ":50055"
	operation = "divide"
)

// OperationServer implementa o serviço OperationService
type OperationServer struct {
	pb.UnimplementedOperationServiceServer
	operation string
}

// NewOperationServer cria um novo servidor de operação
func NewOperationServer(op string) *OperationServer {
	return &OperationServer{
		operation: op,
	}
}

// Execute executa a operação
func (s *OperationServer) Execute(ctx context.Context, req *pb.OperationRequest) (*pb.OperationResponse, error) {
	// Extrai clientID do expressionID
	clientID := "UNKNOWN"
	parts := strings.Split(req.ExpressionId, "_expr_")
	if len(parts) > 0 {
		clientID = parts[0]
	}

	serverName := strings.ToUpper(s.operation)
	log.Printf("[%s] [%s] Recebida operação: %s(%v) [Step: %s]", serverName, clientID, req.Operation, req.Numbers, req.StepId)

	// Valida operação
	if req.Operation != s.operation {
		log.Printf("[%s] [%s] Operação inválida: esperado %s, recebido %s", serverName, clientID, s.operation, req.Operation)
		return &pb.OperationResponse{
			ExpressionId: req.ExpressionId,
			StepId:       req.StepId,
			Error: &pb.ErrorInfo{
				Code:    "INVALID_OPERATION",
				Message: "Este servidor só processa operações " + s.operation,
			},
		}, nil
	}

	// Executa operação
	result, err := grpcOps.ExecuteOperation(req.Operation, req.Numbers)
	if err != nil {
		log.Printf("[%s] [%s] Erro ao executar operação: %v", serverName, clientID, err)
		errorCode := "EXECUTION_ERROR"
		if err.Error() == "divisão por zero" {
			errorCode = "DIV_BY_ZERO"
		}
		return &pb.OperationResponse{
			ExpressionId: req.ExpressionId,
			StepId:       req.StepId,
			Error: &pb.ErrorInfo{
				Code:    errorCode,
				Message: err.Error(),
			},
		}, nil
	}

	log.Printf("[%s] [%s] Operação executada com sucesso: %f", serverName, clientID, result)
	return &pb.OperationResponse{
		ExpressionId: req.ExpressionId,
		StepId:       req.StepId,
		Result:       result,
	}, nil
}

func main() {
	serverName := strings.ToUpper(operation)
	log.Printf("[%s] Iniciando servidor de operação: %s", serverName, operation)

	// Cria listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[%s] Falha ao criar listener: %v", serverName, err)
	}

	// Cria servidor gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterOperationServiceServer(grpcServer, NewOperationServer(operation))

	log.Printf("[%s] Servidor %s escutando em %s", serverName, operation, port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}

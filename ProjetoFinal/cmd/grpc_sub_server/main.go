package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Monterazo/Atividades-IF711/ProjetoFinal/proto"
	grpcOps "github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/grpc"
	"google.golang.org/grpc"
)

const (
	port = ":50053"
	operation = "subtract"
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
	log.Printf("Recebida operação: %s(%v) [Step: %s]", req.Operation, req.Numbers, req.StepId)

	// Valida operação
	if req.Operation != s.operation {
		log.Printf("Operação inválida: esperado %s, recebido %s", s.operation, req.Operation)
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
		log.Printf("Erro ao executar operação: %v", err)
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

	log.Printf("Operação executada com sucesso: %f", result)
	return &pb.OperationResponse{
		ExpressionId: req.ExpressionId,
		StepId:       req.StepId,
		Result:       result,
	}, nil
}

func main() {
	log.Printf("Iniciando servidor de operação: %s", operation)

	// Cria listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Falha ao criar listener: %v", err)
	}

	// Cria servidor gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterOperationServiceServer(grpcServer, NewOperationServer(operation))

	log.Printf("Servidor %s escutando em %s", operation, port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}

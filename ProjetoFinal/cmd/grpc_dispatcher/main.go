package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Monterazo/Atividades-IF711/ProjetoFinal/proto"
	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/core"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// DispatcherServer implementa o serviço CalculatorService
type DispatcherServer struct {
	pb.UnimplementedCalculatorServiceServer
	parser        *core.Parser
	serverAddrs   map[string]string
	serverClients map[string]pb.OperationServiceClient
}

// NewDispatcherServer cria um novo servidor dispatcher
func NewDispatcherServer() *DispatcherServer {
	serverAddrs := map[string]string{
		"add":      "localhost:50052",
		"subtract": "localhost:50053",
		"multiply": "localhost:50054",
		"divide":   "localhost:50055",
	}

	return &DispatcherServer{
		parser:        core.NewParser(),
		serverAddrs:   serverAddrs,
		serverClients: make(map[string]pb.OperationServiceClient),
	}
}

// connectToServers estabelece conexões com os servidores de operação
func (s *DispatcherServer) connectToServers() error {
	for operation, addr := range s.serverAddrs {
		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
		if err != nil {
			return fmt.Errorf("falha ao conectar ao servidor %s: %v", operation, err)
		}
		s.serverClients[operation] = pb.NewOperationServiceClient(conn)
		log.Printf("Conectado ao servidor %s em %s", operation, addr)
	}
	return nil
}

// Calculate processa uma expressão matemática
func (s *DispatcherServer) Calculate(ctx context.Context, req *pb.ExpressionRequest) (*pb.ExpressionResponse, error) {
	log.Printf("Recebida expressão: %s (ID: %s)", req.Expression, req.ExpressionId)

	// Parse da expressão
	steps, err := s.parser.Parse(req.Expression)
	if err != nil {
		log.Printf("Erro ao fazer parse da expressão: %v", err)
		return &pb.ExpressionResponse{
			ExpressionId: req.ExpressionId,
			Error: &pb.ErrorInfo{
				Code:    "PARSE_ERROR",
				Message: fmt.Sprintf("Erro ao fazer parse da expressão: %v", err),
			},
		}, nil
	}

	log.Printf("Expressão parseada em %d steps", len(steps))

	// Executa cada step
	results := make(map[string]float64)
	for i, step := range steps {
		// Substitui referências a resultados anteriores
		numbers := make([]float64, len(step.Numbers))
		copy(numbers, step.Numbers)

		// Cria contexto com timeout
		deadline := time.Duration(req.DeadlineMs) * time.Millisecond
		stepCtx, cancel := context.WithTimeout(ctx, deadline)
		defer cancel()

		// Prepara requisição de operação
		opReq := &pb.OperationRequest{
			ExpressionId: req.ExpressionId,
			StepId:       fmt.Sprintf("%s_step%d", req.ExpressionId, i),
			Operation:    step.Operation,
			Numbers:      numbers,
			DeadlineMs:   req.DeadlineMs,
		}

		log.Printf("Executando step %d: %s(%v)", i, step.Operation, numbers)

		// Obtém cliente do servidor de operação
		client, ok := s.serverClients[step.Operation]
		if !ok {
			return &pb.ExpressionResponse{
				ExpressionId: req.ExpressionId,
				Error: &pb.ErrorInfo{
					Code:    "UNKNOWN_OPERATION",
					Message: fmt.Sprintf("Operação desconhecida: %s", step.Operation),
				},
			}, nil
		}

		// Executa operação
		opResp, err := client.Execute(stepCtx, opReq)
		if err != nil {
			log.Printf("Erro ao executar step %d: %v", i, err)
			return &pb.ExpressionResponse{
				ExpressionId: req.ExpressionId,
				Error: &pb.ErrorInfo{
					Code:    "EXECUTION_ERROR",
					Message: fmt.Sprintf("Erro ao executar operação: %v", err),
				},
			}, nil
		}

		if opResp.Error != nil {
			log.Printf("Erro retornado pelo servidor: %s - %s", opResp.Error.Code, opResp.Error.Message)
			return &pb.ExpressionResponse{
				ExpressionId: req.ExpressionId,
				Error:        opResp.Error,
			}, nil
		}

		results[opReq.StepId] = opResp.Result
		log.Printf("Step %d completado: resultado = %f", i, opResp.Result)

		// Atualiza números para próximos steps se necessário
		if i < len(steps)-1 {
			// O resultado atual será usado no próximo step
			if i+1 < len(steps) {
				steps[i+1].Numbers[0] = opResp.Result
			}
		}

		// Se for o último step, retorna o resultado
		if i == len(steps)-1 {
			log.Printf("Expressão calculada com sucesso: %f", opResp.Result)
			return &pb.ExpressionResponse{
				ExpressionId: req.ExpressionId,
				Result:       opResp.Result,
			}, nil
		}
	}

	// Fallback (não deveria chegar aqui)
	return &pb.ExpressionResponse{
		ExpressionId: req.ExpressionId,
		Error: &pb.ErrorInfo{
			Code:    "INTERNAL_ERROR",
			Message: "Erro interno ao processar expressão",
		},
	}, nil
}

func main() {
	log.Println("Iniciando Dispatcher gRPC...")

	// Cria o servidor
	server := NewDispatcherServer()

	// Aguarda um pouco para os servidores de operação iniciarem
	log.Println("Aguardando servidores de operação...")
	time.Sleep(2 * time.Second)

	// Conecta aos servidores de operação
	if err := server.connectToServers(); err != nil {
		log.Fatalf("Erro ao conectar aos servidores: %v", err)
	}

	// Cria listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Falha ao criar listener: %v", err)
	}

	// Cria servidor gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, server)

	log.Printf("Dispatcher escutando em %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir: %v", err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	pb "github.com/Monterazo/Atividades-IF711/ProjetoFinal/proto"
	"google.golang.org/grpc"
)

type BenchmarkResult struct {
	ClientID       int
	RequestID      int
	Duration       time.Duration
	Success        bool
	ErrorMessage   string
	ExpressionID   string
}

type Statistics struct {
	TotalRequests     int
	SuccessfulReqs    int
	FailedReqs        int
	TotalDuration     time.Duration
	AverageLatency    time.Duration
	MinLatency        time.Duration
	MaxLatency        time.Duration
	P50Latency        time.Duration
	P95Latency        time.Duration
	P99Latency        time.Duration
	Throughput        float64
}

var (
	expression      = flag.String("expr", "((4+3)*2)/5", "Expressão matemática para testar")
	numClients      = flag.Int("clients", 10, "Número de clientes simultâneos")
	reqsPerClient   = flag.Int("reqs", 100, "Número de requisições por cliente")
	dispatcherAddr  = flag.String("dispatcher", "localhost:50051", "Endereço do dispatcher")
	timeoutMs       = flag.Int("timeout", 30000, "Timeout em milissegundos")
	verbose         = flag.Bool("v", false, "Modo verboso (mostra cada requisição)")
)

func main() {
	flag.Parse()

	log.Printf("🚀 Benchmark gRPC - Calculadora Distribuída")
	log.Printf("===========================================")
	log.Printf("Expressão: %s", *expression)
	log.Printf("Clientes simultâneos: %d", *numClients)
	log.Printf("Requisições por cliente: %d", *reqsPerClient)
	log.Printf("Total de requisições: %d", (*numClients) * (*reqsPerClient))
	log.Printf("Dispatcher: %s", *dispatcherAddr)
	log.Printf("Timeout: %dms", *timeoutMs)
	log.Printf("===========================================\n")

	// Canal para coletar resultados
	results := make(chan BenchmarkResult, (*numClients) * (*reqsPerClient))
	var wg sync.WaitGroup

	// Marca início do benchmark
	startTime := time.Now()

	// Lança clientes concorrentes
	for clientID := 0; clientID < *numClients; clientID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runClient(id, *expression, *reqsPerClient, results)
		}(clientID)
	}

	// Aguarda conclusão de todos os clientes
	wg.Wait()
	close(results)

	totalDuration := time.Since(startTime)

	// Processa e exibe resultados
	stats := calculateStatistics(results, totalDuration)
	displayResults(stats)
}

func runClient(clientID int, expr string, numReqs int, results chan<- BenchmarkResult) {
	// Conecta ao dispatcher
	conn, err := grpc.Dial(*dispatcherAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Printf("❌ [Cliente %d] Falha ao conectar: %v", clientID, err)
		for i := 0; i < numReqs; i++ {
			results <- BenchmarkResult{
				ClientID:     clientID,
				RequestID:    i,
				Success:      false,
				ErrorMessage: fmt.Sprintf("Falha ao conectar: %v", err),
			}
		}
		return
	}
	defer conn.Close()

	client := pb.NewCalculatorServiceClient(conn)

	if *verbose {
		log.Printf("✅ [Cliente %d] Conectado ao dispatcher", clientID)
	}

	// Envia requisições
	for reqID := 0; reqID < numReqs; reqID++ {
		expressionID := fmt.Sprintf("CLIENT-%d-REQ-%d-%d", clientID, reqID, time.Now().UnixNano())

		req := &pb.ExpressionRequest{
			ExpressionId: expressionID,
			Expression:   expr,
			DeadlineMs:   int64(*timeoutMs),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutMs)*time.Millisecond)

		startReq := time.Now()
		resp, err := client.Calculate(ctx, req)
		duration := time.Since(startReq)
		cancel()

		result := BenchmarkResult{
			ClientID:     clientID,
			RequestID:    reqID,
			Duration:     duration,
			ExpressionID: expressionID,
		}

		if err != nil {
			result.Success = false
			result.ErrorMessage = err.Error()
			if *verbose {
				log.Printf("❌ [Cliente %d | Req %d] Erro: %v (tempo: %v)", clientID, reqID, err, duration)
			}
		} else if resp.Error != nil {
			result.Success = false
			result.ErrorMessage = fmt.Sprintf("[%s] %s", resp.Error.Code, resp.Error.Message)
			if *verbose {
				log.Printf("❌ [Cliente %d | Req %d] Erro: %s (tempo: %v)", clientID, reqID, result.ErrorMessage, duration)
			}
		} else {
			result.Success = true
			if *verbose {
				log.Printf("✅ [Cliente %d | Req %d] Resultado: %f (tempo: %v)", clientID, reqID, resp.Result, duration)
			}
		}

		results <- result
	}

	if *verbose {
		log.Printf("🏁 [Cliente %d] Finalizou todas as requisições", clientID)
	}
}

func calculateStatistics(resultsChan <-chan BenchmarkResult, totalDuration time.Duration) Statistics {
	var durations []time.Duration
	stats := Statistics{}

	for result := range resultsChan {
		stats.TotalRequests++
		if result.Success {
			stats.SuccessfulReqs++
			durations = append(durations, result.Duration)
		} else {
			stats.FailedReqs++
			if *verbose {
				log.Printf("Erro na requisição %s: %s", result.ExpressionID, result.ErrorMessage)
			}
		}
	}

	if len(durations) == 0 {
		return stats
	}

	// Ordena durações para calcular percentis
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	// Calcula estatísticas
	stats.TotalDuration = totalDuration
	stats.MinLatency = durations[0]
	stats.MaxLatency = durations[len(durations)-1]

	var sum time.Duration
	for _, d := range durations {
		sum += d
	}
	stats.AverageLatency = sum / time.Duration(len(durations))

	// Percentis
	stats.P50Latency = durations[len(durations)*50/100]
	stats.P95Latency = durations[int(math.Min(float64(len(durations)*95/100), float64(len(durations)-1)))]
	stats.P99Latency = durations[int(math.Min(float64(len(durations)*99/100), float64(len(durations)-1)))]

	// Throughput (requisições por segundo)
	stats.Throughput = float64(stats.SuccessfulReqs) / totalDuration.Seconds()

	return stats
}

func displayResults(stats Statistics) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 RESULTADOS DO BENCHMARK")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\n📈 Requisições:\n")
	fmt.Printf("   Total:        %d\n", stats.TotalRequests)
	fmt.Printf("   Sucesso:      %d (%.2f%%)\n", stats.SuccessfulReqs,
		float64(stats.SuccessfulReqs)*100/float64(stats.TotalRequests))
	fmt.Printf("   Falhas:       %d (%.2f%%)\n", stats.FailedReqs,
		float64(stats.FailedReqs)*100/float64(stats.TotalRequests))

	fmt.Printf("\n⏱️  Latência:\n")
	fmt.Printf("   Mínima:       %v\n", stats.MinLatency)
	fmt.Printf("   Média:        %v\n", stats.AverageLatency)
	fmt.Printf("   Máxima:       %v\n", stats.MaxLatency)
	fmt.Printf("   P50:          %v\n", stats.P50Latency)
	fmt.Printf("   P95:          %v\n", stats.P95Latency)
	fmt.Printf("   P99:          %v\n", stats.P99Latency)

	fmt.Printf("\n🚀 Desempenho:\n")
	fmt.Printf("   Duração total:    %v\n", stats.TotalDuration)
	fmt.Printf("   Throughput:       %.2f req/s\n", stats.Throughput)

	fmt.Println("\n" + strings.Repeat("=", 60))
}

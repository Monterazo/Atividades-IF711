package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Monterazo/Atividades-IF711/ProjetoFinal/internal/rabbitmq"
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
	rabbitmqURL     = flag.String("url", "amqp://guest:guest@localhost:5672/", "URL do RabbitMQ")
	timeoutMs       = flag.Int("timeout", 30000, "Timeout em milissegundos")
	verbose         = flag.Bool("v", false, "Modo verboso (mostra cada requisição)")
)

func main() {
	flag.Parse()

	log.Printf("🚀 Benchmark RabbitMQ - Calculadora Distribuída")
	log.Printf("===============================================")
	log.Printf("Expressão: %s", *expression)
	log.Printf("Clientes simultâneos: %d", *numClients)
	log.Printf("Requisições por cliente: %d", *reqsPerClient)
	log.Printf("Total de requisições: %d", (*numClients) * (*reqsPerClient))
	log.Printf("RabbitMQ: %s", *rabbitmqURL)
	log.Printf("Timeout: %dms", *timeoutMs)
	log.Printf("===============================================\n")

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
	clientName := fmt.Sprintf("BENCH-CLIENT-%d", clientID)

	// Conecta ao RabbitMQ
	conn, err := rabbitmq.NewConnection(*rabbitmqURL)
	if err != nil {
		log.Printf("❌ [Cliente %d] Erro ao conectar ao RabbitMQ: %v", clientID, err)
		for i := 0; i < numReqs; i++ {
			results <- BenchmarkResult{
				ClientID:     clientID,
				RequestID:    i,
				Success:      false,
				ErrorMessage: fmt.Sprintf("Erro ao conectar: %v", err),
			}
		}
		return
	}
	defer conn.Close()

	// Declara filas necessárias
	if err := conn.DeclareQueue(rabbitmq.RequestQueue); err != nil {
		log.Printf("❌ [Cliente %d] Erro ao declarar fila de requests: %v", clientID, err)
		return
	}
	if err := conn.DeclareQueue(rabbitmq.ResponseQueue); err != nil {
		log.Printf("❌ [Cliente %d] Erro ao declarar fila de responses: %v", clientID, err)
		return
	}

	// Inicia consumidor de respostas
	msgs, err := conn.Consume(rabbitmq.ResponseQueue)
	if err != nil {
		log.Printf("❌ [Cliente %d] Erro ao consumir fila de responses: %v", clientID, err)
		return
	}

	if *verbose {
		log.Printf("✅ [Cliente %d] Conectado ao RabbitMQ", clientID)
	}

	// Canal para respostas deste cliente
	responseChan := make(chan rabbitmq.ExpressionResponse, numReqs)
	pendingRequests := make(map[string]time.Time)
	var mu sync.Mutex

	// Goroutine para processar respostas
	go func() {
		for msg := range msgs {
			var resp rabbitmq.ExpressionResponse
			if err := json.Unmarshal(msg.Body, &resp); err != nil {
				if *verbose {
					log.Printf("❌ [Cliente %d] Erro ao decodificar resposta: %v", clientID, err)
				}
				msg.Nack(false, false)
				continue
			}

			// Verifica se é resposta para este cliente
			if strings.HasPrefix(resp.ExpressionID, clientName) {
				responseChan <- resp
				msg.Ack(false)
			} else {
				// Rejeita mensagem que não é para este cliente
				msg.Nack(false, true)
			}
		}
	}()

	// Goroutine para processar respostas e calcular latência
	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		received := 0
		for received < numReqs {
			select {
			case resp := <-responseChan:
				mu.Lock()
				startTime, exists := pendingRequests[resp.ExpressionID]
				if exists {
					delete(pendingRequests, resp.ExpressionID)
				}
				mu.Unlock()

				if !exists {
					if *verbose {
						log.Printf("⚠️  [Cliente %d] Resposta recebida para requisição desconhecida: %s", clientID, resp.ExpressionID)
					}
					continue
				}

				duration := time.Since(startTime)
				result := BenchmarkResult{
					ClientID:     clientID,
					RequestID:    received,
					Duration:     duration,
					ExpressionID: resp.ExpressionID,
				}

				if resp.Error != nil {
					result.Success = false
					result.ErrorMessage = fmt.Sprintf("[%s] %s", resp.Error.Code, resp.Error.Message)
					if *verbose {
						log.Printf("❌ [Cliente %d | Req %d] Erro: %s (tempo: %v)", clientID, received, result.ErrorMessage, duration)
					}
				} else {
					result.Success = true
					if *verbose {
						log.Printf("✅ [Cliente %d | Req %d] Resultado: %f (tempo: %v)", clientID, received, resp.Result, duration)
					}
				}

				results <- result
				received++

			case <-time.After(time.Duration(*timeoutMs) * time.Millisecond):
				// Timeout: marca requisições pendentes como falhas
				mu.Lock()
				for exprID, startTime := range pendingRequests {
					results <- BenchmarkResult{
						ClientID:     clientID,
						Duration:     time.Since(startTime),
						ExpressionID: exprID,
						Success:      false,
						ErrorMessage: "Timeout aguardando resposta",
					}
					if *verbose {
						log.Printf("❌ [Cliente %d] Timeout para requisição %s", clientID, exprID)
					}
				}
				pendingRequests = make(map[string]time.Time)
				received = numReqs
				mu.Unlock()
			}
		}
	}()

	// Envia todas as requisições
	for reqID := 0; reqID < numReqs; reqID++ {
		expressionID := fmt.Sprintf("%s-REQ-%d-%d", clientName, reqID, time.Now().UnixNano())

		req := rabbitmq.ExpressionRequest{
			ExpressionID: expressionID,
			Expression:   expr,
			DeadlineMs:   int64(*timeoutMs),
		}

		reqBytes, err := json.Marshal(req)
		if err != nil {
			results <- BenchmarkResult{
				ClientID:     clientID,
				RequestID:    reqID,
				Success:      false,
				ErrorMessage: fmt.Sprintf("Erro ao serializar: %v", err),
			}
			continue
		}

		mu.Lock()
		pendingRequests[expressionID] = time.Now()
		mu.Unlock()

		if err := conn.Publish(rabbitmq.RequestQueue, reqBytes); err != nil {
			mu.Lock()
			delete(pendingRequests, expressionID)
			mu.Unlock()

			results <- BenchmarkResult{
				ClientID:     clientID,
				RequestID:    reqID,
				Success:      false,
				ErrorMessage: fmt.Sprintf("Erro ao enviar: %v", err),
			}
			if *verbose {
				log.Printf("❌ [Cliente %d | Req %d] Erro ao enviar: %v", clientID, reqID, err)
			}
		}
	}

	// Aguarda processamento de todas as respostas
	resultsWg.Wait()

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

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// --- Estrutura para armazenar métricas de desempenho ---

type MetricasDesempenho struct {
	tempoInicio        time.Time
	tempoFim           time.Time
	tempoTotalExecucao time.Duration
	tempoPrimeiraOnda  time.Duration
	tempoCancelamentos time.Duration
	tempoSegundaOnda   time.Duration
	mutex              sync.Mutex
}

// --- Estruturas de Comunicação (Canais) ---

// Resposta unificada para requisições de reserva/cancelamento
type ReservaResponse struct {
	sucesso    bool
	mensagem   string
	tipoQuarto string
}

// Requisição de Reserva (inclui canal para resposta)
type ReservaRequest struct {
	clienteID  int
	tipoQuarto string
	response   chan ReservaResponse
}

// Requisição de Cancelamento
type CancelamentoRequest struct {
	clienteID  int
	tipoQuarto string
	response   chan ReservaResponse
}

// Requisição de Estatísticas
type StatsRequest struct {
	response chan map[string]map[string]int32
}

// --- Estruturas de Dados do Hotel ---

// Estrutura que representa um tipo de quarto
type TipoQuarto struct {
	nome               string
	totalQuartos       int32
	quartosDisponiveis int32
	reservasConcluidas int32
	reservasFalharam   int32
}

// Estrutura que representa o sistema de reservas de hotel
type SistemaReservasHotel struct {
	quartos        map[string]*TipoQuarto
	canal          chan string
	wg             sync.WaitGroup
	requestChannel chan interface{}
	metricas       *MetricasDesempenho // NOVO: métricas de desempenho
}

// Cria um novo sistema de reservas de hotel
func NovoSistemaReservasHotel() *SistemaReservasHotel {
	sistema := &SistemaReservasHotel{
		quartos:        make(map[string]*TipoQuarto),
		canal:          make(chan string, 200),
		requestChannel: make(chan interface{}, 50),
		metricas:       &MetricasDesempenho{}, // Inicializa métricas
	}

	sistema.AdicionarTipoQuarto("Standard", 15)
	sistema.AdicionarTipoQuarto("Luxo", 8)
	sistema.AdicionarTipoQuarto("Suite", 5)
	sistema.AdicionarTipoQuarto("Presidencial", 2)

	return sistema
}

// Adiciona um tipo de quarto ao sistema
func (s *SistemaReservasHotel) AdicionarTipoQuarto(nome string, quantidade int32) {
	s.quartos[nome] = &TipoQuarto{
		nome:               nome,
		totalQuartos:       quantidade,
		quartosDisponiveis: quantidade,
		reservasConcluidas: 0,
		reservasFalharam:   0,
	}
}

// --- Métodos de Medição de Desempenho ---

// Inicia a medição de tempo
func (s *SistemaReservasHotel) IniciarMedicao() {
	s.metricas.mutex.Lock()
	defer s.metricas.mutex.Unlock()
	s.metricas.tempoInicio = time.Now()
}

// Finaliza a medição de tempo
func (s *SistemaReservasHotel) FinalizarMedicao() {
	s.metricas.mutex.Lock()
	defer s.metricas.mutex.Unlock()
	s.metricas.tempoFim = time.Now()
	s.metricas.tempoTotalExecucao = s.metricas.tempoFim.Sub(s.metricas.tempoInicio)
}

// Obtém as métricas de desempenho
func (s *SistemaReservasHotel) ObterMetricas() MetricasDesempenho {
	s.metricas.mutex.Lock()
	defer s.metricas.mutex.Unlock()
	return *s.metricas
}

// Exibe relatório de desempenho
func (s *SistemaReservasHotel) ExibirRelatorioDesempenho() {
	metricas := s.ObterMetricas()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("⏱️  MÉTRICAS DE DESEMPENHO - Chanel")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("⏰ Tempo total de execução: %v\n", metricas.tempoTotalExecucao)
	fmt.Printf("📊 Tempo em milissegundos: %.2f ms\n", float64(metricas.tempoTotalExecucao.Microseconds())/1000.0)
	fmt.Printf("📊 Tempo em segundos: %.4f s\n", metricas.tempoTotalExecucao.Seconds())

	if metricas.tempoPrimeiraOnda > 0 {
		fmt.Printf("👥 Tempo primeira onda: %v (%.2f ms)\n",
			metricas.tempoPrimeiraOnda,
			float64(metricas.tempoPrimeiraOnda.Microseconds())/1000.0)
	}
	if metricas.tempoCancelamentos > 0 {
		fmt.Printf("🔄 Tempo cancelamentos: %v (%.2f ms)\n",
			metricas.tempoCancelamentos,
			float64(metricas.tempoCancelamentos.Microseconds())/1000.0)
	}
	if metricas.tempoSegundaOnda > 0 {
		fmt.Printf("👥 Tempo segunda onda: %v (%.2f ms)\n",
			metricas.tempoSegundaOnda,
			float64(metricas.tempoSegundaOnda.Microseconds())/1000.0)
	}
}

// --- Goroutine Monitor (Gerente de Reservas) ---

func (s *SistemaReservasHotel) GerenteReservas() {
	defer fmt.Println("\n⚠️ Gerente de Reservas parou de processar requisições.")
	for req := range s.requestChannel {
		switch r := req.(type) {
		case ReservaRequest:
			s.executarReserva(r)
		case CancelamentoRequest:
			s.executarCancelamento(r)
		case StatsRequest:
			s.executarObterEstatisticas(r)
		default:
			s.enviarMensagem(fmt.Sprintf("ERRO: Requisição desconhecida recebida: %T", req), 0)
		}
	}
}

// Lógica de Reserva
func (s *SistemaReservasHotel) executarReserva(req ReservaRequest) {
	quarto, existe := s.quartos[req.tipoQuarto]

	if !existe {
		mensagem := fmt.Sprintf("❌ Cliente %d: Tipo de quarto '%s' não existe!", req.clienteID, req.tipoQuarto)
		req.response <- ReservaResponse{sucesso: false, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
		return
	}

	if quarto.quartosDisponiveis > 0 {
		quarto.quartosDisponiveis--
		quarto.reservasConcluidas++

		mensagem := fmt.Sprintf("✅ Cliente %d: Quarto %s RESERVADO! Quartos %s restantes: %d/%d",
			req.clienteID, req.tipoQuarto, req.tipoQuarto, quarto.quartosDisponiveis, quarto.totalQuartos)

		req.response <- ReservaResponse{sucesso: true, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
		return
	}

	quarto.reservasFalharam++
	mensagem := fmt.Sprintf("❌ Cliente %d: Quarto %s INDISPONÍVEL - todos ocupados (%d/%d)",
		req.clienteID, req.tipoQuarto, 0, quarto.totalQuartos)

	req.response <- ReservaResponse{sucesso: false, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
}

// Lógica de Cancelamento
func (s *SistemaReservasHotel) executarCancelamento(req CancelamentoRequest) {
	quarto, existe := s.quartos[req.tipoQuarto]

	if !existe {
		mensagem := fmt.Sprintf("❌ Cliente %d: Não é possível cancelar - tipo '%s' não existe!", req.clienteID, req.tipoQuarto)
		req.response <- ReservaResponse{sucesso: false, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
		return
	}

	if quarto.quartosDisponiveis < quarto.totalQuartos {
		quarto.quartosDisponiveis++
	}

	mensagem := fmt.Sprintf("🔄 Cliente %d: Cancelamento de Quarto %s processado. Disponíveis: %d/%d",
		req.clienteID, req.tipoQuarto, quarto.quartosDisponiveis, quarto.totalQuartos)

	req.response <- ReservaResponse{sucesso: true, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
}

// Lógica de Obter Estatísticas
func (s *SistemaReservasHotel) executarObterEstatisticas(req StatsRequest) {
	stats := make(map[string]map[string]int32)

	for nome, quarto := range s.quartos {
		stats[nome] = map[string]int32{
			"disponiveis": quarto.quartosDisponiveis,
			"total":       quarto.totalQuartos,
			"reservados":  quarto.reservasConcluidas,
			"negadas":     quarto.reservasFalharam,
		}
	}
	req.response <- stats
}

// --- Métodos Chamados pelas Goroutines de Cliente ---

func (s *SistemaReservasHotel) TentarReservaQuarto(clienteID int, tipoQuarto string) {
	defer s.wg.Done()

	tempoProcessamento := time.Duration(rand.Intn(150)) * time.Millisecond
	time.Sleep(tempoProcessamento)

	responseChan := make(chan ReservaResponse, 1)

	req := ReservaRequest{
		clienteID:  clienteID,
		tipoQuarto: tipoQuarto,
		response:   responseChan,
	}

	s.requestChannel <- req
	resposta := <-responseChan
	s.enviarMensagem(resposta.mensagem, clienteID)
}

func (s *SistemaReservasHotel) CancelarReservaQuarto(clienteID int, tipoQuarto string) {
	defer s.wg.Done()

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	responseChan := make(chan ReservaResponse, 1)

	req := CancelamentoRequest{
		clienteID:  clienteID,
		tipoQuarto: tipoQuarto,
		response:   responseChan,
	}

	s.requestChannel <- req
	resposta := <-responseChan
	s.enviarMensagem(resposta.mensagem, clienteID)
}

func (s *SistemaReservasHotel) ObterEstatisticasGerais() map[string]map[string]int32 {
	responseChan := make(chan map[string]map[string]int32, 1)
	req := StatsRequest{response: responseChan}
	s.requestChannel <- req
	stats := <-responseChan
	return stats
}

// --- Lógica de Log ---

func (s *SistemaReservasHotel) enviarMensagem(mensagem string, clienteID int) {
	select {
	case s.canal <- mensagem:
	case <-time.After(15 * time.Millisecond):
		fmt.Printf("⚠️ Timeout - mensagem perdida do cliente %d\n", clienteID)
	}
}

func (s *SistemaReservasHotel) ProcessarMensagens() {
	go func() {
		for mensagem := range s.canal {
			fmt.Println(mensagem)
		}
	}()
}

// --- Função Principal ---

func main() {
	fmt.Println("🏨 Sistema de Reservas de Hotel - Canais (Monitor Goroutine)")
	fmt.Println("==================================================")

	hotel := NovoSistemaReservasHotel()

	// INICIA A MEDIÇÃO DE TEMPO
	hotel.IniciarMedicao()

	// Inicia o Gerente de Reservas
	go hotel.GerenteReservas()

	// Inicia o processador de mensagens
	hotel.ProcessarMensagens()

	fmt.Println("🛏️  TIPOS DE QUARTOS DISPONÍVEIS:")
	fmt.Println("   • Standard: 15 quartos")
	fmt.Println("   • Luxo: 8 quartos")
	fmt.Println("   • Suíte: 5 quartos")
	fmt.Println("   • Presidencial: 2 quartos")
	fmt.Println()

	rand.Seed(time.Now().UnixNano())

	tiposQuartos := []string{"Standard", "Luxo", "Suite", "Presidencial"}
	pesosQuartos := []int{50, 25, 15, 10}

	escolherTipoQuarto := func() string {
		total := 0
		for _, peso := range pesosQuartos {
			total += peso
		}

		r := rand.Intn(total)
		soma := 0

		for i, peso := range pesosQuartos {
			soma += peso
			if r < soma {
				return tiposQuartos[i]
			}
		}
		return tiposQuartos[0]
	}

	// Primeira onda de clientes
	totalClientes := 40
	fmt.Printf("👥 Primeira onda: %d clientes fazendo reservas...\n\n", totalClientes)

	inicioPrimeiraOnda := time.Now()
	for i := 1; i <= totalClientes; i++ {
		hotel.wg.Add(1)
		tipoQuarto := escolherTipoQuarto()
		go hotel.TentarReservaQuarto(i, tipoQuarto)
	}
	hotel.wg.Wait()
	hotel.metricas.tempoPrimeiraOnda = time.Since(inicioPrimeiraOnda)

	// Cancelamentos
	inicioCancelamentos := time.Now()
	hotel.wg.Add(1)
	go func() {
		defer hotel.wg.Done()
		time.Sleep(300 * time.Millisecond)
		fmt.Println("\n🔄 Processando cancelamentos...\n")

		cancelamentos := map[string][]int{
			"Standard":     {101, 102},
			"Luxo":         {103},
			"Suite":        {104},
			"Presidencial": {105},
		}

		for tipo, clientes := range cancelamentos {
			for _, clienteID := range clientes {
				hotel.wg.Add(1)
				go hotel.CancelarReservaQuarto(clienteID, tipo)
			}
		}
	}()
	hotel.wg.Wait()
	hotel.metricas.tempoCancelamentos = time.Since(inicioCancelamentos)

	// Segunda onda de clientes
	inicioSegundaOnda := time.Now()
	hotel.wg.Add(1)
	go func() {
		defer hotel.wg.Done()
		time.Sleep(500 * time.Millisecond)
		fmt.Println("\n👥 Segunda onda: novos clientes...\n")

		for i := 201; i <= 210; i++ {
			hotel.wg.Add(1)
			tipoQuarto := escolherTipoQuarto()
			go hotel.TentarReservaQuarto(i, tipoQuarto)
		}
	}()
	hotel.wg.Wait()
	hotel.metricas.tempoSegundaOnda = time.Since(inicioSegundaOnda)

	// Exibe estatísticas
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 ESTATÍSTICAS DETALHADAS POR TIPO DE QUARTO")
	fmt.Println(strings.Repeat("=", 60))

	stats := hotel.ObterEstatisticasGerais()
	totalReservadas := int32(0)
	totalNegadas := int32(0)
	totalDisponiveis := int32(0)
	totalQuartos := int32(0)

	ordem := []string{"Standard", "Luxo", "Suite", "Presidencial"}

	for _, tipo := range ordem {
		if dados, existe := stats[tipo]; existe {
			fmt.Printf("\n🛏️  QUARTO %s:\n", strings.ToUpper(tipo))
			fmt.Printf("   • Disponíveis: %d/%d quartos\n", dados["disponiveis"], dados["total"])
			fmt.Printf("   • Reservas confirmadas: %d\n", dados["reservados"])
			fmt.Printf("   • Reservas negadas: %d\n", dados["negadas"])

			if dados["reservados"]+dados["negadas"] > 0 {
				taxa := float64(dados["reservados"]) / float64(dados["reservados"]+dados["negadas"]) * 100
				fmt.Printf("   • Taxa de sucesso: %.1f%%\n", taxa)
			}

			totalReservadas += dados["reservados"]
			totalNegadas += dados["negadas"]
			totalDisponiveis += dados["disponiveis"]
			totalQuartos += dados["total"]
		}
	}

	// Fecha canais
	close(hotel.requestChannel)
	time.Sleep(100 * time.Millisecond)
	close(hotel.canal)
	time.Sleep(100 * time.Millisecond)

	// FINALIZA A MEDIÇÃO DE TEMPO
	hotel.FinalizarMedicao()

	// Estatísticas gerais
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Printf("🏨 RESUMO GERAL DO HOTEL:\n")
	fmt.Printf("   • Total de quartos: %d\n", totalQuartos)
	fmt.Printf("   • Quartos ocupados: %d\n", totalQuartos-totalDisponiveis)
	fmt.Printf("   • Quartos disponíveis: %d\n", totalDisponiveis)
	fmt.Printf("   • Taxa de ocupação: %.1f%%\n", float64(totalQuartos-totalDisponiveis)/float64(totalQuartos)*100)

	fmt.Printf("\n📈 ESTATÍSTICAS DE RESERVAS:\n")
	fmt.Printf("   • Reservas confirmadas: %d\n", totalReservadas)
	fmt.Printf("   • Reservas negadas: %d\n", totalNegadas)
	fmt.Printf("   • Total de tentativas: %d\n", totalReservadas+totalNegadas)

	if totalReservadas+totalNegadas > 0 {
		taxaGeral := float64(totalReservadas) / float64(totalReservadas+totalNegadas) * 100
		fmt.Printf("   • Taxa de sucesso geral: %.1f%%\n", taxaGeral)
	}

	// EXIBE RELATÓRIO DE DESEMPENHO
	hotel.ExibirRelatorioDesempenho()

	fmt.Println("\n✨ Sistema de reservas finalizado!")
}

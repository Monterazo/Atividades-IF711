package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Estrutura para armazenar métricas de desempenho
type MetricasDesempenho struct {
	tempoInicio        time.Time
	tempoFim           time.Time
	tempoTotalExecucao time.Duration
	tempoPrimeiraOnda  time.Duration
	tempoCancelamentos time.Duration
	tempoSegundaOnda   time.Duration
	tempoProcessamento time.Duration
	mutex              sync.Mutex
}

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
	quartos  map[string]*TipoQuarto
	canal    chan string
	wg       sync.WaitGroup
	mutex    sync.RWMutex
	metricas *MetricasDesempenho
}

// Cria um novo sistema de reservas de hotel
func NovoSistemaReservasHotel() *SistemaReservasHotel {
	sistema := &SistemaReservasHotel{
		quartos:  make(map[string]*TipoQuarto),
		canal:    make(chan string, 200),
		metricas: &MetricasDesempenho{},
	}

	// Inicializa os diferentes tipos de quartos
	sistema.AdicionarTipoQuarto("Standard", 15)
	sistema.AdicionarTipoQuarto("Luxo", 8)
	sistema.AdicionarTipoQuarto("Suite", 5)
	sistema.AdicionarTipoQuarto("Presidencial", 2)

	return sistema
}

// Adiciona um tipo de quarto ao sistema
func (s *SistemaReservasHotel) AdicionarTipoQuarto(nome string, quantidade int32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.quartos[nome] = &TipoQuarto{
		nome:               nome,
		totalQuartos:       quantidade,
		quartosDisponiveis: quantidade,
		reservasConcluidas: 0,
		reservasFalharam:   0,
	}
}

// Método para tentar fazer uma reserva de quarto
func (s *SistemaReservasHotel) TentarReservaQuarto(clienteID int, tipoQuarto string) {
	defer s.wg.Done()

	// Simula tempo de processamento variável
	tempoProcessamento := time.Duration(rand.Intn(150)) * time.Millisecond
	time.Sleep(tempoProcessamento)

	s.mutex.RLock()
	quarto, existe := s.quartos[tipoQuarto]
	s.mutex.RUnlock()

	if !existe {
		mensagem := fmt.Sprintf("❌ Cliente %d: Tipo de quarto '%s' não existe!", clienteID, tipoQuarto)
		s.enviarMensagem(mensagem, clienteID)
		return
	}

	// Tenta decrementar o contador atomicamente
	quartosAtuais := atomic.LoadInt32(&quarto.quartosDisponiveis)

	if quartosAtuais > 0 {
		// Usa compare-and-swap para garantir atomicidade
		if atomic.CompareAndSwapInt32(&quarto.quartosDisponiveis, quartosAtuais, quartosAtuais-1) {
			// Reserva bem-sucedida
			atomic.AddInt32(&quarto.reservasConcluidas, 1)

			quartosRestantes := atomic.LoadInt32(&quarto.quartosDisponiveis)
			mensagem := fmt.Sprintf("✅ Cliente %d: Quarto %s RESERVADO! Quartos %s restantes: %d/%d",
				clienteID, tipoQuarto, tipoQuarto, quartosRestantes, quarto.totalQuartos)

			s.enviarMensagem(mensagem, clienteID)
			return
		}
	}

	// Reserva falhou - sem quartos disponíveis
	atomic.AddInt32(&quarto.reservasFalharam, 1)

	mensagem := fmt.Sprintf("❌ Cliente %d: Quarto %s INDISPONÍVEL - todos ocupados (%d/%d)",
		clienteID, tipoQuarto, 0, quarto.totalQuartos)

	s.enviarMensagem(mensagem, clienteID)
}

// Método para cancelar uma reserva (liberar quarto)
func (s *SistemaReservasHotel) CancelarReservaQuarto(clienteID int, tipoQuarto string) {
	defer s.wg.Done()

	// Simula tempo de processamento
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	s.mutex.RLock()
	quarto, existe := s.quartos[tipoQuarto]
	s.mutex.RUnlock()

	if !existe {
		mensagem := fmt.Sprintf("❌ Cliente %d: Não é possível cancelar - tipo '%s' não existe!", clienteID, tipoQuarto)
		s.enviarMensagem(mensagem, clienteID)
		return
	}

	// Incrementa atomicamente o número de quartos disponíveis
	novosQuartos := atomic.AddInt32(&quarto.quartosDisponiveis, 1)

	// Garante que não ultrapasse o total
	if novosQuartos > quarto.totalQuartos {
		atomic.StoreInt32(&quarto.quartosDisponiveis, quarto.totalQuartos)
		novosQuartos = quarto.totalQuartos
	}

	mensagem := fmt.Sprintf("🔄 Cliente %d: Cancelamento de Quarto %s processado. Disponíveis: %d/%d",
		clienteID, tipoQuarto, novosQuartos, quarto.totalQuartos)

	s.enviarMensagem(mensagem, clienteID)
}

// Método auxiliar para enviar mensagens com timeout
func (s *SistemaReservasHotel) enviarMensagem(mensagem string, clienteID int) {
	select {
	case s.canal <- mensagem:
	case <-time.After(15 * time.Millisecond):
		fmt.Printf("⚠️ Timeout - mensagem perdida do cliente %d\n", clienteID)
	}
}

// Goroutine para processar mensagens do canal
func (s *SistemaReservasHotel) ProcessarMensagens() {
	go func() {
		for mensagem := range s.canal {
			fmt.Println(mensagem)
		}
	}()
}

// Método para obter estatísticas de um tipo de quarto
func (s *SistemaReservasHotel) ObterEstatisticasQuarto(tipoQuarto string) (int32, int32, int32, int32) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if quarto, existe := s.quartos[tipoQuarto]; existe {
		return atomic.LoadInt32(&quarto.quartosDisponiveis),
			quarto.totalQuartos,
			atomic.LoadInt32(&quarto.reservasConcluidas),
			atomic.LoadInt32(&quarto.reservasFalharam)
	}
	return 0, 0, 0, 0
}

// Método para obter estatísticas gerais
func (s *SistemaReservasHotel) ObterEstatisticasGerais() map[string]map[string]int32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := make(map[string]map[string]int32)

	for nome, quarto := range s.quartos {
		stats[nome] = map[string]int32{
			"disponiveis": atomic.LoadInt32(&quarto.quartosDisponiveis),
			"total":       quarto.totalQuartos,
			"reservados":  atomic.LoadInt32(&quarto.reservasConcluidas),
			"negadas":     atomic.LoadInt32(&quarto.reservasFalharam),
		}
	}

	return stats
}

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
	fmt.Println("⏱️  MÉTRICAS DE DESEMPENHO - Mutex")
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

func main() {
	fmt.Println("🏨 Sistema de Reservas de Hotel - Tipos de Quartos")
	fmt.Println("==================================================")

	// Cria o sistema de reservas
	hotel := NovoSistemaReservasHotel()

	// INICIA A MEDIÇÃO DE TEMPO
	hotel.IniciarMedicao()

	// Inicia o processador de mensagens
	hotel.ProcessarMensagens()

	// Exibe configuração inicial
	fmt.Println("🛏️  TIPOS DE QUARTOS DISPONÍVEIS:")
	fmt.Println("   • Standard: 15 quartos")
	fmt.Println("   • Luxo: 8 quartos")
	fmt.Println("   • Suíte: 5 quartos")
	fmt.Println("   • Presidencial: 2 quartos")
	fmt.Println()

	// Seed para números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Tipos de quartos disponíveis
	tiposQuartos := []string{"Standard", "Luxo", "Suite", "Presidencial"}

	// Pesos para simular preferência (Standard mais procurado)
	pesosQuartos := []int{50, 25, 15, 10}

	// Função para escolher tipo de quarto baseado em peso
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

	// Fecha o canal após delay
	time.Sleep(150 * time.Millisecond)
	close(hotel.canal)
	time.Sleep(100 * time.Millisecond)

	// FINALIZA A MEDIÇÃO DE TEMPO
	hotel.FinalizarMedicao()

	// Exibe estatísticas detalhadas
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

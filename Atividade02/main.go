package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

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
	// Canal de resposta para receber o resultado da operação
	response chan ReservaResponse
}

// Requisição de Cancelamento
type CancelamentoRequest struct {
	clienteID  int
	tipoQuarto string
	response   chan ReservaResponse
}

// Requisição de Estatísticas
type StatsRequest struct {
	// Canal de resposta para receber o mapa de estatísticas
	response chan map[string]map[string]int32
}

// --- Estruturas de Dados do Hotel ---

// Estrutura que representa um tipo de quarto
type TipoQuarto struct {
	nome               string
	totalQuartos       int32
	quartosDisponiveis int32 // Gerenciado apenas pelo GerenteReservas
	reservasConcluidas int32 // Gerenciado apenas pelo GerenteReservas
	reservasFalharam   int32 // Gerenciado apenas pelo GerenteReservas
	// NENHUM mutex ou atomic aqui, pois o acesso é serializado pelo canal.
}

// Estrutura que representa o sistema de reservas de hotel
type SistemaReservasHotel struct {
	quartos map[string]*TipoQuarto // Acessado apenas pelo GerenteReservas
	canal   chan string            // Canal para comunicação de logs/mensagens (mantido da Atividade 1)
	wg      sync.WaitGroup
	// NOVO: Canal principal para todas as requisições de estado (Reserva, Cancelamento, Stats)
	requestChannel chan interface{}
}

// Cria um novo sistema de reservas de hotel
func NovoSistemaReservasHotel() *SistemaReservasHotel {
	sistema := &SistemaReservasHotel{
		quartos:        make(map[string]*TipoQuarto),
		canal:          make(chan string, 200),
		requestChannel: make(chan interface{}, 50), // Canal de requisições
	}

	// Inicializa os diferentes tipos de quartos (sem mutex, pois está na inicialização single-threaded)
	sistema.AdicionarTipoQuarto("Standard", 15)
	sistema.AdicionarTipoQuarto("Luxo", 8)
	sistema.AdicionarTipoQuarto("Suite", 5)
	sistema.AdicionarTipoQuarto("Presidencial", 2)

	return sistema
}

// Adiciona um tipo de quarto ao sistema (usado apenas na inicialização)
func (s *SistemaReservasHotel) AdicionarTipoQuarto(nome string, quantidade int32) {
	s.quartos[nome] = &TipoQuarto{
		nome:               nome,
		totalQuartos:       quantidade,
		quartosDisponiveis: quantidade,
		reservasConcluidas: 0,
		reservasFalharam:   0,
	}
}

// --- Goroutine Monitor (Gerente de Reservas) ---

// Goroutine central que gerencia todo o estado do hotel (o inventário de quartos).
// Ela é a única que lê e escreve no mapa 's.quartos'.
func (s *SistemaReservasHotel) GerenteReservas() {
	defer fmt.Println("\n⚠️ Gerente de Reservas parou de processar requisições.")
	for req := range s.requestChannel {
		// O switch type garante que apenas uma requisição seja processada por vez,
		// eliminando a necessidade de mutex ou atomic para o estado.
		switch r := req.(type) {
		case ReservaRequest:
			s.executarReserva(r)
		case CancelamentoRequest:
			s.executarCancelamento(r)
		case StatsRequest:
			s.executarObterEstatisticas(r)
		default:
			// Log de erro para tipo de requisição desconhecido
			s.enviarMensagem(fmt.Sprintf("ERRO: Requisição desconhecida recebida: %T", req), 0)
		}
	}
}

// Lógica de Reserva (executada APENAS pelo GerenteReservas)
func (s *SistemaReservasHotel) executarReserva(req ReservaRequest) {
	quarto, existe := s.quartos[req.tipoQuarto]

	if !existe {
		mensagem := fmt.Sprintf("❌ Cliente %d: Tipo de quarto '%s' não existe!", req.clienteID, req.tipoQuarto)
		req.response <- ReservaResponse{sucesso: false, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
		return
	}

	if quarto.quartosDisponiveis > 0 {
		// Modificação de estado segura, pois está serializada pelo canal
		quarto.quartosDisponiveis--
		quarto.reservasConcluidas++

		mensagem := fmt.Sprintf("✅ Cliente %d: Quarto %s RESERVADO! Quartos %s restantes: %d/%d",
			req.clienteID, req.tipoQuarto, req.tipoQuarto, quarto.quartosDisponiveis, quarto.totalQuartos)

		req.response <- ReservaResponse{sucesso: true, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
		return
	}

	// Reserva falhou - sem quartos disponíveis
	quarto.reservasFalharam++
	mensagem := fmt.Sprintf("❌ Cliente %d: Quarto %s INDISPONÍVEL - todos ocupados (%d/%d)",
		req.clienteID, req.tipoQuarto, 0, quarto.totalQuartos)

	req.response <- ReservaResponse{sucesso: false, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
}

// Lógica de Cancelamento (executada APENAS pelo GerenteReservas)
func (s *SistemaReservasHotel) executarCancelamento(req CancelamentoRequest) {
	quarto, existe := s.quartos[req.tipoQuarto]

	if !existe {
		mensagem := fmt.Sprintf("❌ Cliente %d: Não é possível cancelar - tipo '%s' não existe!", req.clienteID, req.tipoQuarto)
		req.response <- ReservaResponse{sucesso: false, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
		return
	}

	// Libera o quarto se o número for menor que o total
	if quarto.quartosDisponiveis < quarto.totalQuartos {
		quarto.quartosDisponiveis++
	}

	mensagem := fmt.Sprintf("🔄 Cliente %d: Cancelamento de Quarto %s processado. Disponíveis: %d/%d",
		req.clienteID, req.tipoQuarto, quarto.quartosDisponiveis, quarto.totalQuartos)

	req.response <- ReservaResponse{sucesso: true, mensagem: mensagem, tipoQuarto: req.tipoQuarto}
}

// Lógica de Obter Estatísticas (executada APENAS pelo GerenteReservas)
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
	// Envia o resultado de volta para o canal de resposta
	req.response <- stats
}

// --- Métodos Chamados pelas Goroutines de Cliente (Fachada) ---

// Método para tentar fazer uma reserva de quarto (envia requisição ao GerenteReservas)
func (s *SistemaReservasHotel) TentarReservaQuarto(clienteID int, tipoQuarto string) {
	defer s.wg.Done()

	// Simula tempo de processamento variável
	tempoProcessamento := time.Duration(rand.Intn(150)) * time.Millisecond
	time.Sleep(tempoProcessamento)

	// Cria canal de resposta exclusivo para esta requisição
	responseChan := make(chan ReservaResponse, 1)

	// Cria a requisição
	req := ReservaRequest{
		clienteID:  clienteID,
		tipoQuarto: tipoQuarto,
		response:   responseChan,
	}

	// Envia a requisição para o gerente
	s.requestChannel <- req

	// Aguarda a resposta do gerente
	resposta := <-responseChan

	// Envia a mensagem de log
	s.enviarMensagem(resposta.mensagem, clienteID)
}

// Método para cancelar uma reserva (liberar quarto)
func (s *SistemaReservasHotel) CancelarReservaQuarto(clienteID int, tipoQuarto string) {
	defer s.wg.Done()

	// Simula tempo de processamento
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

// Método para obter estatísticas gerais (envia requisição ao GerenteReservas)
func (s *SistemaReservasHotel) ObterEstatisticasGerais() map[string]map[string]int32 {
	responseChan := make(chan map[string]map[string]int32, 1)

	req := StatsRequest{response: responseChan}

	// Envia a requisição de leitura ao gerente
	s.requestChannel <- req

	// Aguarda o resultado
	stats := <-responseChan
	return stats
}

// --- Lógica de Log (Mantida da Atividade 1) ---

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

// --- Função Principal (Main) ---

func main() {
	fmt.Println("🏨 Sistema de Reservas de Hotel - Canais (Monitor Goroutine)")
	fmt.Println("==================================================")

	// Cria o sistema de reservas
	hotel := NovoSistemaReservasHotel()

	// Inicia o Gerente de Reservas (a goroutine monitora)
	go hotel.GerenteReservas()

	// Inicia o processador de mensagens
	hotel.ProcessarMensagens()

	// Exibe configuração inicial
	fmt.Println("🛏️ 	TIPOS DE QUARTOS DISPONÍVEIS:")
	fmt.Println(" 	• Standard: 15 quartos")
	fmt.Println(" 	• Luxo: 8 quartos")
	fmt.Println(" 	• Suíte: 5 quartos")
	fmt.Println(" 	• Presidencial: 2 quartos")
	fmt.Println()

	// Seed para números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Tipos de quartos disponíveis
	tiposQuartos := []string{"Standard", "Luxo", "Suite", "Presidencial"}

	// Pesos para simular preferência (Standard mais procurado)
	pesosQuartos := []int{50, 25, 15, 10} // Porcentagem de preferência

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

	for i := 1; i <= totalClientes; i++ {
		hotel.wg.Add(1)
		tipoQuarto := escolherTipoQuarto()
		go hotel.TentarReservaQuarto(i, tipoQuarto)
	}

	// Aguarda um pouco e faz alguns cancelamentos
	hotel.wg.Add(1)
	go func() {
		defer hotel.wg.Done()
		time.Sleep(300 * time.Millisecond)
		fmt.Println("\n🔄 Processando cancelamentos...\n")

		// Cancelamentos de diferentes tipos
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

	// Segunda onda de clientes aproveitando cancelamentos
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

	// Aguarda todas as goroutines de clientes e operações secundárias terminarem
	hotel.wg.Wait()

	// Exibe estatísticas detalhadas
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 ESTATÍSTICAS DETALHADAS POR TIPO DE QUARTO")
	fmt.Println(strings.Repeat("=", 60))

	stats := hotel.ObterEstatisticasGerais() // Esta chamada também usa o canal para leitura
	totalReservadas := int32(0)
	totalNegadas := int32(0)

	// Fecha o canal de requisições do gerente
	close(hotel.requestChannel)
	time.Sleep(100 * time.Millisecond)

	// Fecha o canal de logs após delay
	close(hotel.canal)
	time.Sleep(100 * time.Millisecond)
	totalDisponiveis := int32(0)
	totalQuartos := int32(0)

	// Ordena tipos por importância
	ordem := []string{"Standard", "Luxo", "Suite", "Presidencial"}

	for _, tipo := range ordem {
		if dados, existe := stats[tipo]; existe {
			fmt.Printf("\n🛏️ 	QUARTO %s:\n", strings.ToUpper(tipo))
			fmt.Printf(" 	• Disponíveis: %d/%d quartos\n", dados["disponiveis"], dados["total"])
			fmt.Printf(" 	• Reservas confirmadas: %d\n", dados["reservados"])
			fmt.Printf(" 	• Reservas negadas: %d\n", dados["negadas"])

			if dados["reservados"]+dados["negadas"] > 0 {
				taxa := float64(dados["reservados"]) / float64(dados["reservados"]+dados["negadas"]) * 100
				fmt.Printf(" 	• Taxa de sucesso: %.1f%%\n", taxa)
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
	fmt.Printf(" 	• Total de quartos: %d\n", totalQuartos)
	fmt.Printf(" 	• Quartos ocupados: %d\n", totalQuartos-totalDisponiveis)
	fmt.Printf(" 	• Quartos disponíveis: %d\n", totalDisponiveis)
	fmt.Printf(" 	• Taxa de ocupação: %.1f%%\n", float64(totalQuartos-totalDisponiveis)/float64(totalQuartos)*100)

	fmt.Printf("\n📈 ESTATÍSTICAS DE RESERVAS:\n")
	fmt.Printf(" 	• Reservas confirmadas: %d\n", totalReservadas)
	fmt.Printf(" 	• Reservas negadas: %d\n", totalNegadas)
	fmt.Printf(" 	• Total de tentativas: %d\n", totalReservadas+totalNegadas)

	if totalReservadas+totalNegadas > 0 {
		taxaGeral := float64(totalReservadas) / float64(totalReservadas+totalNegadas) * 100
		fmt.Printf(" 	• Taxa de sucesso geral: %.1f%%\n", taxaGeral)
	}

	fmt.Println("\n✨ Sistema de reservas finalizado!")
}

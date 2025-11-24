package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Filas de requisição e resposta
	RequestQueue  = "calculator.requests"
	ResponseQueue = "calculator.responses"

	// Filas de operações
	AddQueue      = "operations.add"
	SubtractQueue = "operations.subtract"
	MultiplyQueue = "operations.multiply"
	DivideQueue   = "operations.divide"
	ResultsQueue  = "operations.results"
)

// Connection encapsula uma conexão RabbitMQ
type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewConnection cria uma nova conexão com o RabbitMQ
func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("falha ao criar canal: %v", err)
	}

	return &Connection{
		conn:    conn,
		channel: channel,
	}, nil
}

// Close fecha a conexão
func (c *Connection) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// DeclareQueue declara uma fila
func (c *Connection) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,  // nome
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// Publish publica uma mensagem em uma fila
func (c *Connection) Publish(queue string, body []byte) error {
	return c.channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
}

// Consume consome mensagens de uma fila
func (c *Connection) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

// SetupQueues declara todas as filas necessárias
func SetupQueues(conn *Connection) error {
	queues := []string{
		RequestQueue,
		ResponseQueue,
		AddQueue,
		SubtractQueue,
		MultiplyQueue,
		DivideQueue,
		ResultsQueue,
	}

	for _, queue := range queues {
		if err := conn.DeclareQueue(queue); err != nil {
			return fmt.Errorf("falha ao declarar fila %s: %v", queue, err)
		}
		log.Printf("Fila declarada: %s", queue)
	}

	return nil
}

// GetQueueForOperation retorna o nome da fila para uma operação
func GetQueueForOperation(operation string) string {
	switch operation {
	case "add":
		return AddQueue
	case "subtract":
		return SubtractQueue
	case "multiply":
		return MultiplyQueue
	case "divide":
		return DivideQueue
	default:
		return ""
	}
}

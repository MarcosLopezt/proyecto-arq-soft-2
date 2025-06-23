package queues

import (
	"encoding/json"
	"fmt"
	"log"
	cursos "search_cursos/domain"

	"github.com/streadway/amqp"
)

type RabbitConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	QueueName string
}

type Rabbit struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewRabbit(config RabbitConfig) Rabbit {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		log.Fatalf("Error al obtener conexión Rabbit: %v", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Error al crear canal Rabbit: %v", err)
	}
	queue, err := channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al declarar la cola Rabbit: %v", err)
	}
	return Rabbit{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}
}

func (queue Rabbit) StartConsumer(handler func(cursos.CursoNew)) error {
	messages, err := queue.channel.Consume(
		queue.queue.Name,
		"",
		true, 
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error al registrar consumidor: %w", err)
	}

	go func() {
		for msg := range messages {
			var cursoUpdate cursos.CursoNew
			if err := json.Unmarshal(msg.Body, &cursoUpdate); err != nil {
				log.Printf("Error al deserializar mensaje: %v", err)
				continue
			}
			handler(cursoUpdate)
		}
	}()

	return nil
}

func (queue Rabbit) Close() {
	if err := queue.channel.Close(); err != nil {
		log.Printf("Error al cerrar canal Rabbit: %v", err)
	}
	if err := queue.connection.Close(); err != nil {
		log.Printf("Error al cerrar conexión Rabbit: %v", err)
	}
}

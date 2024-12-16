package queues

import (
	cursos "cursos/models"
	"encoding/json"
	"fmt"
	"log"

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

func NewRabbit(config RabbitConfig) (*Rabbit, error) {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		return nil, fmt.Errorf("Error al obtener conexión Rabbit: %w", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("Error al crear canal Rabbit: %w", err)
	}
	queue, err := channel.QueueDeclare(config.QueueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("Error al declarar la cola Rabbit: %w", err)
	}
	return &Rabbit{
		connection: connection,
		channel:    channel,
		queue:      queue,
	}, nil
}

func (r *Rabbit) StartConsumer(handler func(cursos.CursoNew)) error {
	messages, err := r.channel.Consume(
		r.queue.Name,
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

func (r *Rabbit) Close() {
	if err := r.channel.Close(); err != nil {
		log.Printf("Error al cerrar canal Rabbit: %v", err)
	}
	if err := r.connection.Close(); err != nil {
		log.Printf("Error al cerrar conexión Rabbit: %v", err)
	}
}

package queues

import (
	"encoding/json"
	"fmt"
	"log"
	cursos "search_cursos/domain"
	"time"

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
    config     RabbitConfig
}

func NewRabbit(config RabbitConfig) Rabbit {
    var conn *amqp.Connection
    var err error
    for i := 0; i < 5; i++ {
        conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port))
        if err == nil {
            break
        }
        log.Printf("Error conectando a RabbitMQ (intento %d): %v", i+1, err)
        time.Sleep(5 * time.Second)
    }
    if err != nil {
        log.Fatalf("Error al obtener conexión Rabbit: %v", err)
    }
    channel, err := conn.Channel()
    if err != nil {
        log.Fatalf("Error al crear canal Rabbit: %v", err)
    }
    queue, err := channel.QueueDeclare(config.QueueName, true, false, false, false, nil) // Durable queue
    if err != nil {
        log.Fatalf("Error al declarar la cola Rabbit: %v", err)
    }
    return Rabbit{
        connection: conn,
        channel:    channel,
        queue:      queue,
        config:     config,
    }
}

func (queue Rabbit) StartConsumer(handler func(cursos.CursoNew)) error {
    messages, err := queue.channel.Consume(
        queue.queue.Name,
        "",
        false, // Manual acknowledgment
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
			log.Printf("Recibido mensaje: %s", string(msg.Body))
            var cursoUpdate cursos.CursoNew
            if err := json.Unmarshal(msg.Body, &cursoUpdate); err != nil {
                log.Printf("Error al deserializar mensaje: %v", err)
                msg.Nack(false, true) // Requeue message
                continue
            }
			log.Printf("Mensaje deserializado correctamente: %+v", cursoUpdate)
            handler(cursoUpdate)
			log.Printf("Mensaje procesado con éxito, ack")
            msg.Ack(false) // Acknowledge message
        }
    }()

    return nil
}

func (queue Rabbit) Close() {
    if queue.channel != nil {
        if err := queue.channel.Close(); err != nil {
            log.Printf("Error al cerrar canal Rabbit: %v", err)
        }
    }
    if queue.connection != nil {
        if err := queue.connection.Close(); err != nil {
            log.Printf("Error al cerrar conexión Rabbit: %v", err)
        }
    }
}
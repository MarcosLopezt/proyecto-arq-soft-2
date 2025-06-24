package main

import (
	"context"
	"cursos/queues"
	"cursos/router"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var rabbit *queues.Rabbit

func initMongoClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/arq-soft")

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	fmt.Println("Conectado a MongoDB")
	return client, nil
}

func main(){
	var err error
	mongoClient, err = initMongoClient()
	if err != nil {
		log.Fatalf("Error inicializando el cliente de MongoDB: %v", err)
	}
	
	rabbitConfig := queues.RabbitConfig{
		Host:      "localhost", 
		Port:      "5672",
		Username:  "root",
		Password:  "password",
		QueueName: "courses-news",
	}

	rabbit, err = queues.NewRabbit(rabbitConfig)
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		MaxAge:           12 * time.Hour,
	}))

	router.SetupRouter(engine, mongoClient, rabbit)

	if err := engine.Run(":8083"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
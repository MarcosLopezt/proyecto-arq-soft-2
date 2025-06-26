package main

import (
	"context"
	"cursos/queues"
	"cursos/router"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var rabbit *queues.Rabbit

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initMongoClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Usar variables de entorno con valores por defecto para Docker
	mongoHost := getEnv("MONGO_HOST", "mongo")
	mongoPort := getEnv("MONGO_PORT", "27017")
	mongoDB := getEnv("MONGO_DB", "arq-soft")
	
	mongoURI := fmt.Sprintf("mongodb://%s:%s/%s", mongoHost, mongoPort, mongoDB)
	
	clientOptions := options.Client().ApplyURI(mongoURI)

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
		Host:      "rabbitmq", 
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
package main

import (
	"log"
	"search_cursos/controllers"
	"search_cursos/queues"
	"search_cursos/repositories"
	"search_cursos/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Solr
	solrRepo := repositories.NewSolr(repositories.SolrConfig{
		Host:       "solr",    
		Port:       "8983",    
		Collection: "courses",
	})

	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "password",
		QueueName: "courses-news", 
	})

	coursesAPI := repositories.NewHTTP(repositories.HTTPConfig{
		Host: "backend_courses",
		Port: "8083",
	})

	service := services.NewService(solrRepo, coursesAPI)

	controller := controllers.NewController(service)

	if err := eventsQueue.StartConsumer(service.HandleCursoNew); err != nil {
		log.Fatalf("Error running consumer: %v", err)
	}

	router := gin.Default()
	router.GET("/search", controller.Search)
	if err := router.Run(":8085"); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

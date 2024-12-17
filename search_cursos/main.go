package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"search_cursos/controllers"
	"search_cursos/dao"
	"search_cursos/queues"
	"search_cursos/repositories"
	"search_cursos/services"

	"github.com/gin-gonic/gin"
)

func initializeSolr(service services.Service) {
	// Hacer solicitud GET a la API de cursos para obtener todos los cursos
	resp, err := http.Get("http://backend_courses:8083/cursos/all")
	if err != nil {
		log.Fatalf("Error al obtener cursos de la API: %v", err)
	}
	defer resp.Body.Close()

	// Decodificar la respuesta JSON
	var cursos []dao.Curso
	if err := json.NewDecoder(resp.Body).Decode(&cursos); err != nil {
		log.Fatalf("Error al decodificar los cursos: %v", err)
	}

	// Indexar los cursos en Solr
	for _, curso := range cursos {
		cursoDAO := dao.Curso{
			ID:          curso.ID,
			CourseName:  curso.CourseName,
			Description: curso.Description,
			Category:    curso.Category,
			Length:      curso.Length,
		}

		// Indexar curso en Solr
		if _, err := service.Repository.Index(context.Background(), cursoDAO); err != nil {
			log.Printf("Error indexando curso: %v", err)
		} else {
			fmt.Println("Curso indexado correctamente:", curso.ID)
		}
	}
}

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

	initializeSolr(service)
	
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

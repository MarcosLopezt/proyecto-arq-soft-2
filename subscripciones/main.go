package main

import (
	"log"
	"subscripciones/db"
	"subscripciones/router"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){
	err := db.Connect()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	engine := gin.Default()

	// Configuración de CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permitir solicitudes de cualquier origen
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.SetupRouter(engine)

	// Iniciar el servidor en el puerto 8084
	err = engine.Run(":8084")
	if err != nil {
		log.Fatalf("Error al ejecutar el servidor: %v", err)
	}
}
package main

import (
	"log"
	"os"
	"time"
	"users/db"
	"users/router"
	"users/services/cache"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){
	err := db.Connect()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	 // Inicializamos la caché
	 config := cache.MemcachedConfig{
		Host:     os.Getenv("CACHE_HOST"),
		Port:     os.Getenv("CACHE_PORT"),
		Duration: 10 * time.Minute,
	}
    cacheInstance, err := cache.NewCache(config)
	if err != nil {
		log.Fatalf("Error al inicializar la caché: %v", err)
	}
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Permitir solicitudes de cualquier origen
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	router.SetupRouter(engine, cacheInstance)

	err = engine.Run(":8082")
    if err != nil {
        log.Fatalf("Error al ejecutar el servidor: %v", err)
    }
	
}
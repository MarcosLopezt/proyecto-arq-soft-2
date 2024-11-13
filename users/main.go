package main

import (
	"log"
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
	 cacheConfig := cache.CacheConfig{
        MaxSize:      1024 * 1024 * 10, 
        ItemsToPrune: 100,
        Duration:     10 * time.Minute,  // Duración de la caché
    }
    cache := cache.NewCache(cacheConfig)

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Permitir solicitudes de cualquier origen
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	router.SetupRouter(engine, cache)

	err = engine.Run(":8082")
    if err != nil {
        log.Fatalf("Error al ejecutar el servidor: %v", err)
    }
	
}
package main

import (
	"backend/db"
	"backend/router"
	"backend/services/cache"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Creamos una instancia de gin
	err := db.Connect()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Inicializamos la caché
	cacheConfig := cache.CacheConfig{
		MaxSize:      1024 * 1024 * 10, // 10MB, por ejemplo
		ItemsToPrune: 100,
		Duration:     10 * time.Minute,  // Duración de la caché
	}
	cache := cache.NewCache(cacheConfig)

	engine := gin.Default()

	// Configuracion del middleware CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permitir solicitudes de cualquier origen
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Mapeamos las rutas definidas en el router
	router.SetupRouter(engine, cache)

	// Ejecutamos el servidor en el puerto 8080
	err = engine.Run(":8080")
	if err != nil {
		panic(err) 
	}
}

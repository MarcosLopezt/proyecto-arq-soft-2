package test

import (
	"context"
	"fmt"
	"log"
	"time"
	models "users/models"

	"users/services/cache"
)

func main() {
	// Configuración para Memcached
	config := cache.MemcachedConfig{
		Host:     "localhost",       // Cambia por la dirección de tu servidor Memcached
		Port:     "11211",           // Cambia si usas otro puerto
		Duration: 5 * time.Minute,   // Tiempo de expiración de las claves
	}

	// Crear la instancia de la caché
	cacheInstance, err := cache.NewCache(config)
	if err != nil {
		log.Fatalf("Error al inicializar la caché: %v", err)
	}

	// Contexto para operaciones
	ctx := context.Background()

	// Crear un usuario de prueba
	testUser := models.User{
		ID:    1,
		Email: "test@example.com",
	}

	// Guardar el usuario en la caché usando el ID como clave
	_, err = cacheInstance.Create(ctx, testUser)
	if err != nil {
		log.Fatalf("Error al guardar el usuario en caché: %v", err)
	}
	fmt.Println("Usuario guardado exitosamente en la caché.")

	// Intentar recuperar el usuario desde la caché
	retrievedUser, err := cacheInstance.GetUserByID(ctx, fmt.Sprintf("%d", testUser.ID))
	if err != nil {
		log.Fatalf("Error al recuperar el usuario desde la caché: %v", err)
	}
	fmt.Printf("Usuario recuperado desde la caché: %+v\n", retrievedUser)

	// Guardar el usuario en la caché usando el email como clave
	err = cacheInstance.CreateUserByEmail(ctx, &testUser)
	if err != nil {
		log.Fatalf("Error al guardar el usuario en caché por email: %v", err)
	}
	fmt.Println("Usuario guardado exitosamente en la caché con clave de email.")

	// Intentar recuperar el usuario desde la caché por email
	retrievedUserByEmail, err := cacheInstance.GetUserByEmail(ctx, testUser.Email)
	if err != nil {
		log.Fatalf("Error al recuperar el usuario desde la caché por email: %v", err)
	}
	fmt.Printf("Usuario recuperado desde la caché por email: %+v\n", retrievedUserByEmail)
}

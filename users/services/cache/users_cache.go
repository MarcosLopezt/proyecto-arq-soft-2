package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	usersDAO "users/models"

	"github.com/bradfitz/gomemcache/memcache"
)

const (
	keyFormatByID    = "user:%s"
	keyFormatByEmail = "user:email:%s"
)

type MemcachedConfig struct {
	Host     string        // Dirección del host Memcached
	Port     string        // Puerto de Memcached
	Duration time.Duration // Tiempo de expiración para las claves
}


type Cache struct {
	client   *memcache.Client
	duration time.Duration
}


// NewCache inicializa una conexión con el servidor Memcached y verifica la conexión inicial
func NewCache(config MemcachedConfig) (Cache, error) {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// Crear el cliente Memcached
	client := memcache.New(address)

	// Verificar la conexión inicial
	testKey := "test_connection_key"
	testValue := "test_value"
	err := client.Set(&memcache.Item{Key: testKey, Value: []byte(testValue)})
	if err != nil {
		return Cache{}, fmt.Errorf("error al conectar con Memcached en %s: %w", address, err)
	}

	_, err = client.Get(testKey)
	if err != nil {
		return Cache{}, fmt.Errorf("error al verificar Memcached en %s: %w", address, err)
	}

	log.Printf("Conexión a Memcached exitosa en %s", address)

	return Cache{
		client:   client,
		duration: config.Duration,
	}, nil
}

// GetUserByID obtiene un usuario de la caché por ID
func (repository Cache) GetUserByID(ctx context.Context, id string) (usersDAO.User, error) {
	key := fmt.Sprintf(keyFormatByID, id)
	item, err := repository.client.Get(key)
	if err == memcache.ErrCacheMiss {
		fmt.Println("1")
		return usersDAO.User{}, fmt.Errorf("no se encontró el ítem con la clave %s", key)
	} else if err != nil {
		fmt.Println("2")
		return usersDAO.User{}, fmt.Errorf("error obteniendo el ítem con la clave %s: %w", key, err)
	}

	var user usersDAO.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		fmt.Println("3")
		return usersDAO.User{}, fmt.Errorf("error al deserializar el ítem con la clave %s: %w", key, err)
	}

	return user, nil
}

// GetUserByEmail obtiene un usuario de la caché por email
func (repository Cache) GetUserByEmail(ctx context.Context, email string) (usersDAO.User, error) {
	key := fmt.Sprintf(keyFormatByEmail, email)
	item, err := repository.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return usersDAO.User{}, fmt.Errorf("no se encontró el ítem con la clave %s", key)
	} else if err != nil {
		return usersDAO.User{}, fmt.Errorf("error obteniendo el ítem con la clave %s: %w", key, err)
	}

	var user usersDAO.User
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return usersDAO.User{}, fmt.Errorf("error al deserializar el ítem con la clave %s: %w", key, err)
	}

	return user, nil
}

// CreateUserByEmail guarda un usuario en la caché usando el email como clave
func (repository Cache) CreateUserByEmail(ctx context.Context, user *usersDAO.User) error {
	key := fmt.Sprintf(keyFormatByEmail, user.Email)
	fmt.Println("Guardando usuario en caché con clave:", key)

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("error al serializar usuario: %w", err)
	}

	err = repository.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(repository.duration.Seconds()),
	})
	if err != nil {
		return fmt.Errorf("error al guardar usuario en caché: %w", err)
	}

	return nil
}

// Create guarda un usuario en la caché usando el ID como clave
func (repository Cache) Create(ctx context.Context, user usersDAO.User) (string, error) {
	key := fmt.Sprintf(keyFormatByID, fmt.Sprintf("%d", user.ID))
	fmt.Println("Guardando usuario en caché con clave:", key)

	data, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("error al serializar usuario: %w", err)
	}

	err = repository.client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(repository.duration.Seconds()),
	})
	if err != nil {
		return "", fmt.Errorf("error al guardar usuario en caché: %w", err)
	}

	return fmt.Sprintf("%d", user.ID), nil
}

// Invalidate elimina un ítem de la caché
func (repository Cache) Invalidate(ctx context.Context, key string) error {
	err := repository.client.Delete(key)
	if err == memcache.ErrCacheMiss {
		fmt.Printf("No se encontró el ítem con la clave %s para eliminar\n", key)
		return nil
	} else if err != nil {
		return fmt.Errorf("error eliminando el ítem con la clave %s: %w", key, err)
	}
	return nil
}

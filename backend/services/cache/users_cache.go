package cache

import (
	"backend/models/users"
	usersDAO "backend/models/users" // Cambia esto según la estructura de tu proyecto
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/karlseguin/ccache"
)

const (
	keyFormatByID    = "user:%d"
	keyFormatByEmail = "user:email:%s"
)


type CacheConfig struct {
	MaxSize      int64
	ItemsToPrune uint32
	Duration     time.Duration
}

type Cache struct {
	client   *ccache.Cache
	duration time.Duration
}

func NewCache(config CacheConfig) Cache {
	client := ccache.New(ccache.Configure().
		MaxSize(config.MaxSize).
		ItemsToPrune(config.ItemsToPrune))
	return Cache{
		client:   client,
		duration: config.Duration,
	}
}

func (repository Cache) GetUserByID(ctx context.Context, id string) (usersDAO.User, error) {
	key := fmt.Sprintf(keyFormatByID, id)
	item := repository.client.Get(key)
	fmt.Println(key)
	if item == nil {
		return usersDAO.User{}, fmt.Errorf("not found item with key %s", key)
	}
	if item.Expired() {
		return usersDAO.User{}, fmt.Errorf("item with key %s is expired", key)
	}
	userDAO, ok := item.Value().(usersDAO.User)
	if !ok {
		return usersDAO.User{}, fmt.Errorf("error converting item with key %s", key)
	}
	return userDAO, nil
}

func (repository Cache) GetUserByEmail(ctx context.Context, email string) (usersDAO.User, error) {
    key := fmt.Sprintf(keyFormatByEmail, email)
    item := repository.client.Get(key)
    fmt.Println(key)
    
    if item == nil {
        return usersDAO.User{}, fmt.Errorf("no se encontró el ítem con la clave %s", key)
    }
    if item.Expired() {
        return usersDAO.User{}, fmt.Errorf("el ítem con la clave %s ha expirado", key)
    }

    // Cambiar aquí: convertir a []byte antes de deserializar
    userData, ok := item.Value().([]byte)
    if !ok {
        return usersDAO.User{}, fmt.Errorf("error convirtiendo el ítem con la clave %s a []byte", key)
    }

    // Deserializar el usuario desde JSON
    var user usersDAO.User
    if err := json.Unmarshal(userData, &user); err != nil {
        return usersDAO.User{}, fmt.Errorf("error al deserializar el ítem con la clave %s: %w", key, err)
    }

    return user, nil
}


func (repository Cache) CreateUserByEmail(ctx context.Context, user *users.User) error {
	key := fmt.Sprintf(keyFormatByEmail, user.Email)
	fmt.Println("Guardando usuario en caché con clave:", key)

	data, err := json.Marshal(user)
	if err != nil {
        return fmt.Errorf("error al serializar usuario: %w", err)
    }
	repository.client.Set(key, data, repository.duration)
	return nil
}

func (repository Cache) Create(ctx context.Context, user usersDAO.User) (string, error) {
	key := fmt.Sprintf(keyFormatByID, user.ID)
	fmt.Println("saving with duration", repository.duration)
	repository.client.Set(key, user, repository.duration)
	return fmt.Sprintf("%d", user.ID), nil
}

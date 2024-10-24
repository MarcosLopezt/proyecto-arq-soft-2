package dao

import (
	"backend/db"
	"backend/models/users"
	"backend/services/cache"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

//crear nuevo usuario en db

func CreateUser(user *users.User) error {
    return db.DB.Create(user).Error
}

//traer todos los users

func GetAllUsers() ([]users.User, error) {
    var users []users.User
    if err := db.DB.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

//traer user por id

func GetUserByID(ctx context.Context, cache cache.Cache, id uint) (*users.User, error) {
    // primero intenta obtener el usuario de la cache
    
    user, err := cache.GetUserByID(ctx, fmt.Sprintf("%d", id))
    if err != nil {
        return &user, nil
    }

    var userDB users.User  // el user de mysql

    if err := db.DB.First(&userDB, id).Error; err != nil {
        return nil, err
    }

    _ , _= cache.Create(ctx, userDB)

    return &user, nil
}

func GetUserByEmail(ctx context.Context, cache cache.Cache, email string) (*users.User, error) {
    user, err := cache.GetUserByEmail(ctx, email)
    if err == nil {
        fmt.Println("Usuario encontrado en la cache")
        return &user, nil
    }
    var userDB users.User
    if err := db.DB.Where("email = ?", email).First(&userDB).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            fmt.Println("No se encontro el mail en la BD")
            return nil, nil
        }
        return nil, err
    }

    cache.CreateUserByEmail(ctx, &userDB)
    return &userDB, nil
}
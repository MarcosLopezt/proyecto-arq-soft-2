package users_service

import (
	"backend/auth"
	"backend/dao"
	"backend/models/users"
	"backend/services/cache"
	"context"
	"errors"
	"log"

	//"net/http"
	"strconv"

	//"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(cache cache.Cache, request users.LoginRequest) (users.LoginResponse, error) {
	user, err := dao.GetUserByEmail(context.Background(), cache ,request.Email)
	if err != nil {
		return users.LoginResponse{}, err
	}

	if user == nil {
        return users.LoginResponse{}, errors.New("credenciales inválidas")
    }

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return users.LoginResponse{}, errors.New("credenciales inválidas")
	}

	token, err := auth.GenerateAuthToken(user.ID)
	if err != nil {
		return users.LoginResponse{}, err
	}

	return users.LoginResponse{Token: token, Role: user.Role, ID: user.ID}, nil
}

func CreateUser(request users.CreateUserRequest) (users.UserResponse, error) {
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return users.UserResponse{}, err
	}

	user := &users.User{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Role:         request.Role,
	}

	
	if err := dao.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return users.UserResponse{}, err
	}

	
	return users.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func GetUserByID(cache cache.Cache, id string) (users.UserResponse, error) {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return users.UserResponse{}, errors.New("ID inválido")
	}
	user, err := dao.GetUserByID(context.Background(), cache, uint(uid))
	if err != nil {
		return users.UserResponse{}, err
	}

	return users.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

/*

func GetAllUsers() ([]users.UserResponse, error) {
	users, err := dao.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []users.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, users.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return userResponses, nil
}
*/

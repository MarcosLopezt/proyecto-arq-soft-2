package users

import (
	usersDomain "backend/models/users"
	"backend/services/cache"
	usersService "backend/services/users_service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Db *gorm.DB

func LoginHandler(c *gin.Context, cache cache.Cache) {
	var loginRequest usersDomain.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := usersService.Login(cache, loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	c.SetCookie(
		"session_token",
		response.Token,
		int(time.Until(expirationTime).Seconds()),
		"/",   // Path
		"",    // Domain
		false, // Secure
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, response)
}

func CreateUser(c *gin.Context) {
	var createUserRequest usersDomain.CreateUserRequest

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Inicio creacion1")

	user, err := usersService.CreateUser(createUserRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func GetUserByID(c *gin.Context, cache cache.Cache) {
	id := c.Param("id")
	user, err := usersService.GetUserByID(cache, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

/*
func GetAllUsers(c *gin.Context) {
	users, err := usersService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

*/

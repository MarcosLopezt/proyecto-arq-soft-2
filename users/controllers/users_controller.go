package users

import (
	"net/http"
	"os"
	"time"
	usersDomain "users/models"
	usersService "users/services"
	"users/services/cache"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Ping(c *gin.Context) {
	hostname, _ := os.Hostname()
	c.JSON(200, gin.H{
		"message": "Ping recibido",
		"host":    hostname,
	})
}

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

func CreateUser(c *gin.Context, cache cache.Cache) {
	var createUserRequest usersDomain.CreateUserRequest

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := usersService.CreateUser(cache, createUserRequest)
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

func GetInstances(c *gin.Context){
	instances, err := usersService.GetInstances()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"instances": instances,
	})
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

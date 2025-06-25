package router

import (
	controllers "users/controllers"
	"users/services/cache"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r* gin.Engine, cache cache.Cache)*gin.Engine{

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "users-service",
			"timestamp": gin.H{
				"current": "2024-01-01T00:00:00Z",
			},
		})
	})

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/login", func(c *gin.Context){
			controllers.LoginHandler(c, cache)
		})

		userRoutes.POST("register", func(c *gin.Context){
			controllers.CreateUser(c, cache)
		})
		userRoutes.GET("/:id", func(c* gin.Context){
			controllers.GetUserByID(c, cache)
		})

		userRoutes.GET("microservicios", controllers.GetInstances)
	}

	return r
}
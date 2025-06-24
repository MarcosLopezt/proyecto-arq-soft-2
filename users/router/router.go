package router

import (
	controllers "users/controllers"
	"users/services/cache"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r* gin.Engine, cache cache.Cache)*gin.Engine{

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
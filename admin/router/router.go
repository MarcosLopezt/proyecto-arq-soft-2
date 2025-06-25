package router

import (
	"admin/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "admin-service",
			"timestamp": gin.H{
				"current": "2024-01-01T00:00:00Z",
			},
		})
	})

    router.GET("/admin/services", controllers.ListarInstancias)
    router.POST("/admin/services/create", controllers.CrearInstancia)
    router.POST("/admin/services/remove", controllers.EliminarInstancia)
}

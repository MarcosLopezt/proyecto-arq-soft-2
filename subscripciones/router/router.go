package router

import (
	"subscripciones/controllers/comentarios"
	"subscripciones/controllers/files"
	subs "subscripciones/controllers/subs"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "subscriptions-service",
			"timestamp": gin.H{
				"current": "2024-01-01T00:00:00Z",
			},
		})
	})

	 subsRoutes := r.Group("/subscriptions")
	{
	 	subsRoutes.POST("/sub", subs.CreateSubs)
	 	subsRoutes.GET("/get/:user_id", subs.GetSubByUserId)
		subsRoutes.GET("/get/curso/:curso_id", subs.GetSubByCursoId)
	}

	comentRoutes := r.Group("/coments")
	{
		comentRoutes.POST("/coment", comentarios.CreateComent)
		comentRoutes.GET("/:id", comentarios.GetComentsByCourse)
	}

	fileRoutes := r.Group("/files")
	{
		fileRoutes.POST("/upload", files.UploadFile)
		fileRoutes.GET("/file/:curso_id", files.GetFile)
	}

	return r
}
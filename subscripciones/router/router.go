package router

import (
	"subscripciones/controllers/comentarios"
	"subscripciones/controllers/files"
	subs "subscripciones/controllers/subs"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
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
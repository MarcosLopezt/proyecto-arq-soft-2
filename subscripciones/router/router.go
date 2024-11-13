package router

import (
	"subscripciones/controllers/comentarios"
	subs "subscripciones/controllers/subs"
	"subscripciones/controllers/files"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	 subsRoutes := r.Group("/subscriptions")
	{
	 	subsRoutes.POST("/sub", subs.CreateSubs)
	 	subsRoutes.GET("/get/:user_id", subs.GetSubByUserId)
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
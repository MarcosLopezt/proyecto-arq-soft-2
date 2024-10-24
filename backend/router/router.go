package router

import (
	"backend/controllers/comentarios"
	"backend/controllers/cursos"
	"backend/controllers/files"
	"backend/controllers/subscripciones"
	"backend/controllers/users"
	"backend/services/cache"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, cache cache.Cache) *gin.Engine {
	// Grupo de rutas para usuarios
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/login", func(c *gin.Context){
			users.LoginHandler(c, cache)
		})
		userRoutes.POST("/register", users.CreateUser)
		userRoutes.GET("/:id", func(c *gin.Context){
			users.GetUserByID(c, cache)
		})
	}

	// Grupo de rutas para cursos
	courseRoutes := r.Group("/cursos")
	{
		//courseRoutes.Use(middleware.AuthMiddleware())
		courseRoutes.POST("/curso", cursos.CreateCourse)
		courseRoutes.GET("/:course_name", cursos.GetCourseByName)
		courseRoutes.GET("/get/:id", cursos.GetCourseByID)
		courseRoutes.PUT("/update", cursos.UpdateCourse)		
		courseRoutes.DELETE("/delete", cursos.DeleteCourse)
	}

	subsRoutes := r.Group("/subscriptions")
	{
		subsRoutes.POST("/sub", subscripciones.CreateSubs)
		subsRoutes.GET("/get/:user_id", subscripciones.GetSubByUserId)
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

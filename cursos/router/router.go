package router

import (
	controllers "cursos/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(r *gin.Engine, mongoClient *mongo.Client) *gin.Engine {
	courseRoutes := r.Group("/cursos")
	{
		courseRoutes.POST("/curso", func(c *gin.Context) {
			controllers.CreateCourse(c, mongoClient)
		})
		courseRoutes.GET("/:course_name", func(c *gin.Context) {
			controllers.GetCourseByName(c, mongoClient)
		})
		courseRoutes.GET("/get/:id", func(c *gin.Context) {
			controllers.GetCourseByID(c, mongoClient)
		})
		courseRoutes.PUT("/update", func(c *gin.Context) {
			controllers.UpdateCourse(c, mongoClient)
		})
		courseRoutes.DELETE("/delete", func(c *gin.Context) {
			controllers.DeleteCourse(c, mongoClient)
		})
		courseRoutes.GET("/all", func(c *gin.Context){
			controllers.GetAllCourses(c, mongoClient)
		})
	}
	return r
}

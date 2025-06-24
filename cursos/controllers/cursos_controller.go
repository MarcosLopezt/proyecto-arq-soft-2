package cursos

import (
	cursosDomain "cursos/models"
	"cursos/queues"
	cursosService "cursos/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCourse(c *gin.Context, mongoClient *mongo.Client, rabbit *queues.Rabbit) {
    var createCourseRequest cursosDomain.CreateCourseRequest
    if err := c.ShouldBindJSON(&createCourseRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    course, err := cursosService.CreateCourse(mongoClient, createCourseRequest, rabbit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, course)
}
func GetCourseByName(c *gin.Context, mongoClient *mongo.Client) {
	name := c.Param("course_name")

	course, err := cursosService.GetCourseByName(mongoClient, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func GetCourseByID(c *gin.Context, mongoClient *mongo.Client) {
	id := c.Param("id")
	course, err := cursosService.GetCourseByID1(mongoClient, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

func GetAllCourses(c *gin.Context, mongoClient *mongo.Client){
	courses, err := cursosService.GetAllCourses(mongoClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func UpdateCourse(c *gin.Context, mongoClient *mongo.Client, rabbit *queues.Rabbit) {
    var updateCourseRequest cursosDomain.UpdateCourseRequest
    if err := c.ShouldBindJSON(&updateCourseRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    course, err := cursosService.UpdateCourse(mongoClient, updateCourseRequest, rabbit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, course)
}

func DeleteCourse(c *gin.Context, mongoClient *mongo.Client, rabbit *queues.Rabbit) {
	var deleteCourseRequest cursosDomain.DeleteCourseRequest
	if err := c.ShouldBindJSON(&deleteCourseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := cursosService.DeleteCourse(mongoClient, deleteCourseRequest, rabbit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, course)
}

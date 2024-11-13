package cursos

import (
	"context"
	"cursos/dao"
	cursos "cursos/models"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Variable de contexto para MongoDB
var ctx = context.TODO()

// Cambiar la firma de las funciones para que acepten un DAO
func CreateCourse(mongoClient *mongo.Client, request cursos.CreateCourseRequest) (cursos.CreateCourseResponse, error) {
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")

	curso := &cursos.Course{
		CourseName:  request.CourseName,
		Category:    request.Category,
		Length:      request.Length,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := courseDAO.CreateCourse(ctx, curso); err != nil {
		log.Printf("Error creating course: %v", err)
		return cursos.CreateCourseResponse{}, err
	}

	return cursos.CreateCourseResponse{
		Message: "Se creó el curso exitosamente",
	}, nil
}

func GetCourseByName(mongoClient *mongo.Client, courseNameReq string) ([]cursos.GetCourseByNameResponse, error) {
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")

	courses, err := courseDAO.GetCourseByName(ctx, courseNameReq)
	if err != nil {
		return nil, err
	}

	var response []cursos.GetCourseByNameResponse
	for _, course := range courses {
		response = append(response, cursos.GetCourseByNameResponse{
			ID:          course.ID,
			CourseName:  course.CourseName,
			Description: course.Description,
			Length:      course.Length,
			Category:    course.Category,
		})
	}
	return response, nil
}

func GetAllCourses(mongoClient *mongo.Client) ([]cursos.GetCourseByNameResponse, error){
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")

	courses, err := courseDAO.GetAllCourses(ctx)
	if err != nil{
		return nil, err
	}

	var response []cursos.GetCourseByNameResponse
	for _, course := range courses {
		response = append(response, cursos.GetCourseByNameResponse{
			ID:          course.ID,
			CourseName:  course.CourseName,
			Description: course.Description,
			Length:      course.Length,
			Category:    course.Category,
		})
	}

	return response, nil
}

func GetCourseByID(mongoClient *mongo.Client, id string) (cursos.GetCourseByIDResponse, error) {
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")
	curso, err := courseDAO.GetCourseByID(ctx, id)
	if err != nil {
		return cursos.GetCourseByIDResponse{}, err
	}

	if curso == nil {
		return cursos.GetCourseByIDResponse{}, errors.New("curso no encontrado")
	}

	return cursos.GetCourseByIDResponse{
		ID: curso.ID,
		CourseName:  curso.CourseName,
		Category:    curso.Category,
		Description: curso.Description,
		Length:      curso.Length,
	}, nil
}

func UpdateCourse(mongoClient *mongo.Client, request cursos.UpdateCourseRequest) (cursos.UpdateCourseResponse, error) {
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")
	updatedCourse, err := courseDAO.UpdateCourse(ctx, &request)
	if err != nil {
		log.Printf("Error updating course: %v", err)
		return cursos.UpdateCourseResponse{}, err
	}

	return cursos.UpdateCourseResponse{
		ID:          updatedCourse.ID,
		CourseName:  updatedCourse.CourseName,
		Category:    updatedCourse.Category,
		Description: updatedCourse.Description,
		Length:      updatedCourse.Length,
	}, nil
}

func DeleteCourse(mongoClient *mongo.Client, request cursos.DeleteCourseRequest) (cursos.DeleteCourseResponse, error) {
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")
	deleteResponse, err := courseDAO.DeleteCourse(ctx, &request)
	if err != nil {
		log.Printf("Error deleting course: %v", err)
		return cursos.DeleteCourseResponse{}, err
	}

	return cursos.DeleteCourseResponse{
		Message: deleteResponse.Message,
	}, nil
}

package cursos

import (
	"context"
	"cursos/dao"
	cursos "cursos/models"
	"errors"
	"log"

	//"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Variable de contexto para MongoDB
var ctx = context.TODO()

func CreateCourse(mongoClient *mongo.Client, request cursos.CreateCourseRequest) (cursos.CreateCourseResponse, error) {
	courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")

	curso := &cursos.Course{
		CourseName:  request.CourseName,
		Category:    request.Category,
		Length:      request.Length,
		Description: request.Description,
		Cupos: request.Cupos,
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
			Cupos: course.Cupos,
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
			Cupos: course.Cupos,
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
		Cupos: curso.Cupos,
	}, nil
}

type disponibilidadResult struct {
    disponibilidad int
    err            error
}

func GetCourseByID1(mongoClient *mongo.Client, id string) (cursos.GetCourseByIDResponse, error) {
    courseDAO := dao.NewMongoCourseDAO(mongoClient, "arqui_soft", "courses")
    curso, err := courseDAO.GetCourseByID(context.Background(), id)
    if err != nil {
        return cursos.GetCourseByIDResponse{}, err
    }

    if curso == nil {
        return cursos.GetCourseByIDResponse{}, errors.New("curso no encontrado")
    }

    resultCh := make(chan disponibilidadResult, 1) // Canal para disponibilidad o error

    go func() {
        disponibilidad, err := courseDAO.CalcularDisponibilidad(context.Background(), curso.ID)
        resultCh <- disponibilidadResult{disponibilidad: disponibilidad, err: err}
        close(resultCh)
    }()

    result := <-resultCh
    if result.err != nil {
        log.Printf("Error calculando la disponibilidad del curso %s: %v", curso.ID, result.err)
        return cursos.GetCourseByIDResponse{}, result.err
    }

    log.Printf("Disponibilidad actualizada para el curso %s: %d", curso.ID, result.disponibilidad)

    return cursos.GetCourseByIDResponse{
        ID:          curso.ID,
        CourseName:  curso.CourseName,
        Category:    curso.Category,
        Description: curso.Description,
        Length:      curso.Length,
        Cupos:       curso.Cupos,
        Disponibles: result.disponibilidad,
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
		Cupos: updatedCourse.Cupos,
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

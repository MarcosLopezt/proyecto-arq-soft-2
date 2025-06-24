package comentarios_service

import (
	"log"
	comentdao "subscripciones/dao/comentarios"
	"subscripciones/models/comentarios"
)

// Variable global DAO que puede ser reemplazada en tests
var Dao comentdao.ComentarioDAO = comentdao.RealDAO{}

func CreateComent(request comentarios.CreateComentRequest) (comentarios.CreateComentResponse, error) {

	coment := &comentarios.Comentario{
		UserId:  request.UserID,
		CursoId: request.CursoID,
		Texto:   request.Texto,
		Valor:   request.Valor,
	}

	if err := Dao.CreateComent(coment); err != nil {
		log.Printf("Error creating coment: %v", err)
		return comentarios.CreateComentResponse{}, err
	}

	return comentarios.CreateComentResponse{
		Message: "Comentario realizado con exito!",
	}, nil
}

func GetComentByCourse(id int) ([]comentarios.GetComentByCourseResp, error) {
	coments, err := Dao.GetComentByCourse(id)
	if err != nil {
		return nil, err
	}

	var response []comentarios.GetComentByCourseResp
	for _, coment := range coments {
		response = append(response, comentarios.GetComentByCourseResp{
			CursoID: coment.CursoId,
			UserID:  coment.UserId,
			Texto:   coment.Texto,
			Fecha:   coment.Fecha,
			Valor:   coment.Valor,
		})
	}
	return response, nil
}

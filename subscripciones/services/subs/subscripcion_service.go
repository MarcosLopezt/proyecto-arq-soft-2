package subscripcion_service

import (
	"errors"
	"log"
	"strconv"
	dao "subscripciones/dao/subs"
	subscripciones "subscripciones/models/subs"
)

// Declaramos un DAO global que implementa SubsDAO
var Dao dao.SubsDAO = dao.RealDAO{}

func CreateSubs(request subscripciones.CreateSubsRequest) (subscripciones.CreateSubsResponse, error) {
	sub := &subscripciones.Subscription{
		UserID:  request.UserID,
		CourseID: request.CourseID,
	}

	if err := Dao.CreateSubs(sub); err != nil {
		log.Printf("Error creating sub: %v", err)
		return subscripciones.CreateSubsResponse{}, err
	}

	return subscripciones.CreateSubsResponse{
		Message: "Subscripcion realizada con exito!",
	}, nil
}

func GetSubByUserId(id string) ([]subscripciones.GetSubByUserResponse, error) {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return []subscripciones.GetSubByUserResponse{}, errors.New("ID invalido")
	}

	subs, err := Dao.GetSubByUserId(uint(uid))
	if err != nil {
		return []subscripciones.GetSubByUserResponse{}, err
	}

	var response []subscripciones.GetSubByUserResponse
	for _, sub := range subs {
		response = append(response, subscripciones.GetSubByUserResponse{
			ID:       sub.ID,
			UserID:   sub.UserID,
			CourseID: sub.CourseID,
		})
	}

	return response, nil
}

func GetSubByCursoId(id string) (int, error) {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("ID invalido")
	}

	subs, err := Dao.GetSubByCursoId(uint(uid))
	if err != nil {
		return 0, err
	}

	return len(subs), nil
}

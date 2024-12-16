package subscripcion_service

import (
	"errors"
	"log"
	"strconv"
	dao "subscripciones/dao/subs"
	subscripciones "subscripciones/models/subs"
)

func CreateSubs(request subscripciones.CreateSubsRequest) (subscripciones.CreateSubsResponse, error) {

	sub := &subscripciones.Subscription{
		UserID:			 request.UserID,
		CourseID:        request.CourseID,
	}
	
	// cuposResp, err:= dao.GetCupos(request.CourseID)
	// if err != nil {
	// 	log.Printf("Error getting cupos: %v", err)
	// 	return subscripciones.CreateSubsResponse{}, err
	// }

	// if cuposResp.Cupos <= 0 {
	// 	log.Printf("No hay cupos disponibles para el curso %d", request.CourseID)
	// 	return subscripciones.CreateSubsResponse{
	// 		Message: fmt.Sprintf("No se puede suscribir. El curso %d no tiene cupos disponibles.", request.CourseID),
	// 	}, nil
	// }

	if err := dao.CreateSubs(sub); err != nil {
		log.Printf("Error creating sub: %v", err)
		return subscripciones.CreateSubsResponse{}, err
	}

	
	return subscripciones.CreateSubsResponse{
		Message: "Subscripcion realizada con exito!",
	}, nil
}

func GetSubByUserId(id string) ([]subscripciones.GetSubByUserResponse, error){
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil{
		return []subscripciones.GetSubByUserResponse{}, errors.New("ID invalido")
	}

	subs, err := dao.GetSubByUserId(uint(uid))
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

func GetSubByCursoId(id string)(int, error){
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil{
		return 0, errors.New("ID invalido")
	}

	subs, err := dao.GetSubByCursoId(uint(uid))
	if err != nil {
		return 0, err
	}


	return len(subs), nil
}

// func CountSubsByCourse(id string)(int, error){
// 	id, err := strconv.ParseInt(id,10,64)
// 	if err != nil{
// 		return 0, errors.New("ID invalido")
// 	}
// 	subs, err := dao.CountSubsByCourse()
// }


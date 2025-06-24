package comentarios

import (
	"errors"
	"subscripciones/db"
	"subscripciones/models/comentarios"
)

// RealDAO implementa ComentarioDAO
type RealDAO struct{}

func (r RealDAO) CreateComent(coment *comentarios.Comentario) error {
	return db.DB.Create(coment).Error
}

func (r RealDAO) GetComentByCourse(id int) ([]comentarios.Comentario, error) {
	var coments []comentarios.Comentario

	if err := db.DB.Where("curso_id = ?", id).Find(&coments).Error; err != nil {
		return nil, err
	}

	if len(coments) == 0 {
		return nil, errors.New("no hay comentarios en este curso")
	}

	return coments, nil
}

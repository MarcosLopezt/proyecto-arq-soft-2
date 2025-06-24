package comentarios

import "subscripciones/models/comentarios"

type ComentarioDAO interface {
	CreateComent(coment *comentarios.Comentario) error
	GetComentByCourse(id int) ([]comentarios.Comentario, error)
}

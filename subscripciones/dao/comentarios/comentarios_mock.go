package comentarios

import "subscripciones/models/comentarios"

// MockDAO implementa ComentarioDAO para mocks en tests
type MockDAO struct {
	CreateComentFunc      func(*comentarios.Comentario) error
	GetComentByCourseFunc func(int) ([]comentarios.Comentario, error)
}

func (m MockDAO) CreateComent(coment *comentarios.Comentario) error {
	return m.CreateComentFunc(coment)
}

func (m MockDAO) GetComentByCourse(id int) ([]comentarios.Comentario, error) {
	return m.GetComentByCourseFunc(id)
}

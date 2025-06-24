package files

import (
	"subscripciones/models/files"
)

// MockDAO implementa FileDAO para mocks en tests
type MockDAO struct {
	UploadFileFunc func(file *files.File) error
	GetFileFunc    func(cursoId uint) ([]files.File, error)
}

func (m MockDAO) UploadFile(file *files.File) error {
	return m.UploadFileFunc(file)
}

func (m MockDAO) GetFile(cursoId uint) ([]files.File, error) {
	return m.GetFileFunc(cursoId)
}

package files

import (
	"errors"
	"subscripciones/db"
	"subscripciones/models/files"
)

// FileDAO define las operaciones
type FileDAO interface {
	UploadFile(file *files.File) error
	GetFile(cursoId uint) ([]files.File, error)
}

// RealDAO es la implementación real con la base de datos
type RealDAO struct{}

func (r RealDAO) UploadFile(file *files.File) error {
	return db.DB.Create(file).Error
}

func (r RealDAO) GetFile(cursoId uint) ([]files.File, error) {
	var filesList []files.File
	if err := db.DB.Where("curso_id = ?", cursoId).Find(&filesList).Error; err != nil {
		return nil, err
	}

	if len(filesList) == 0 {
		return nil, errors.New("no files found for this course")
	}

	return filesList, nil
}

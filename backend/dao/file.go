package dao

import (
	"backend/db"
	"backend/models/files"
	"errors"
)

func UploadFile(file *files.File) error {
	return db.DB.Create(file).Error
}

func GetFile(cursoId uint) ([]files.File,error){
	var files []files.File
	if err:= db.DB.Where("curso_id = ?", cursoId).Find(&files).Error; err != nil {
		return nil, err
	}

	if len(files) == 0 {
        return nil, errors.New("no files found for this course")
    }

	return files, nil
}
package files_service

import (
	"backend/dao"
	"backend/models/files"
	"errors"
	"log"
	"strconv"
)

func UploadFile(request files.UploadFile) (files.UploadFileResponse, error) {

	file := &files.File{
		CursoID: request.CourseID,
		File: request.File,
	}

	if err := dao.UploadFile(file); err != nil {
		log.Printf("Error uploading file: %v", err)
		return files.UploadFileResponse{}, err
	}


	return files.UploadFileResponse{
		Message: "Archivo subido con exito!",
	}, nil
}


func GetFile(id string) ([]files.GetFileResponse, error){
	uid, err := strconv.ParseUint(id, 10, 32)

	if err != nil{
		return []files.GetFileResponse{}, errors.New("ID invalido")
	}

	filess, err := dao.GetFile(uint(uid))
	if err != nil {
		return nil, err
	}

	var response []files.GetFileResponse
	
    for _, file := range filess {
        response = append(response, files.GetFileResponse{
            ID: file.ID,
			File: file.File,

        })
    }

	return response, nil
}


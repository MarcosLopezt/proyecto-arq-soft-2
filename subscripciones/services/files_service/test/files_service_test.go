package files_service_test

import (
	"errors"
	"subscripciones/dao/files"
	file "subscripciones/models/files"
	"subscripciones/services/files_service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMock() *files.MockDAO {
	mockDAO := &files.MockDAO{}
	files_service.Dao = mockDAO // reemplaza variable global por mock
	return mockDAO
}

func TestUploadFile_Success(t *testing.T) {
	mockDAO := setupMock()

	mockDAO.UploadFileFunc = func(file *file.File) error {
		return nil
	}

	req := file.UploadFile{
		CourseID: 1,
		File:     []byte("test data"),
	}

	resp, err := files_service.UploadFile(req)
	assert.NoError(t, err)
	assert.Equal(t, "Archivo subido con exito!", resp.Message)
}

func TestUploadFile_Error(t *testing.T) {
	mockDAO := setupMock()

	mockDAO.UploadFileFunc = func(file *file.File) error {
		return errors.New("DB error")
	}

	req := file.UploadFile{
		CourseID: 1,
		File:     []byte("test data"),
	}

	_, err := files_service.UploadFile(req)
	assert.Error(t, err)
	assert.Equal(t, "DB error", err.Error())
}

func TestGetFile_Success(t *testing.T) {
	mockDAO := setupMock()

	mockDAO.GetFileFunc = func(cursoId uint) ([]file.File, error) {
		return []file.File{
			{ID: 1, CursoID: cursoId, File: []byte("data")},
		}, nil
	}

	resp, err := files_service.GetFile("1")
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, uint(1), resp[0].ID)
}

func TestGetFile_InvalidID(t *testing.T) {
	setupMock() // solo necesitamos mock para cumplir

	resp, err := files_service.GetFile("abc")
	assert.Error(t, err)
	assert.Equal(t, "ID invalido", err.Error())
	assert.Empty(t, resp)
}

func TestGetFile_DAOError(t *testing.T) {
	mockDAO := setupMock()

	mockDAO.GetFileFunc = func(cursoId uint) ([]file.File, error) {
		return nil, errors.New("DAO error")
	}

	resp, err := files_service.GetFile("1")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "DAO error", err.Error())
}

func TestGetFile_NoResults(t *testing.T) {
	mockDAO := setupMock()

	mockDAO.GetFileFunc = func(cursoId uint) ([]file.File, error) {
		return nil, errors.New("no files found for this course")
	}

	resp, err := files_service.GetFile("1")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "no files found for this course", err.Error())
}

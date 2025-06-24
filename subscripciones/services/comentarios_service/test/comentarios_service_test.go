package comentarios_service_test

import (
	"errors"
	"subscripciones/dao/comentarios"
	comentario "subscripciones/models/comentarios"
	"subscripciones/services/comentarios_service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupMock() comentarios.MockDAO {
	return comentarios.MockDAO{}
}

func TestCreateComent_Success(t *testing.T) {
	mockDAO := setupMock()
	mockDAO.CreateComentFunc = func(coment *comentario.Comentario) error {
		return nil
	}
	comentarios_service.Dao = mockDAO

	req := comentario.CreateComentRequest{
		UserID:  1,
		CursoID: 1,
		Texto:   "Muy bueno",
		Valor:   5,
	}

	resp, err := comentarios_service.CreateComent(req)
	assert.NoError(t, err)
	assert.Equal(t, "Comentario realizado con exito!", resp.Message)
}

func TestCreateComent_Error(t *testing.T) {
	mockDAO := setupMock()
	mockDAO.CreateComentFunc = func(coment *comentario.Comentario) error {
		return errors.New("DB error")
	}
	comentarios_service.Dao = mockDAO

	req := comentario.CreateComentRequest{
		UserID:  1,
		CursoID: 1,
		Texto:   "Muy bueno",
		Valor:   5,
	}

	_, err := comentarios_service.CreateComent(req)
	assert.Error(t, err)
	assert.Equal(t, "DB error", err.Error())
}

func TestGetComentByCourse_Success(t *testing.T) {
	mockDAO := setupMock()
	mockDAO.GetComentByCourseFunc = func(id int) ([]comentario.Comentario, error) {
		return []comentario.Comentario{
			{
				ID:      1,
				UserId:  2,
				CursoId: 3,
				Texto:   "Comentario",
				Fecha:   time.Now(),
				Valor:   4,
			},
		}, nil
	}
	comentarios_service.Dao = mockDAO

	resp, err := comentarios_service.GetComentByCourse(3)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, uint(3), resp[0].CursoID)
}

func TestGetComentByCourse_Error(t *testing.T) {
	mockDAO := setupMock()
	mockDAO.GetComentByCourseFunc = func(id int) ([]comentario.Comentario, error) {
		return nil, errors.New("no hay comentarios en este curso")
	}
	comentarios_service.Dao = mockDAO

	resp, err := comentarios_service.GetComentByCourse(999)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "no hay comentarios en este curso", err.Error())
}

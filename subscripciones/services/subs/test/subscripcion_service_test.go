package subscripcion_service_test

import (
	"errors"
	"subscripciones/dao/subs"
	sub "subscripciones/models/subs"
	service "subscripciones/services/subs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSubs_Success(t *testing.T) {
	mockDao := subs.MockDAO{
		CreateSubsFunc: func(s *sub.Subscription) error {
			return nil
		},
	}
	service.Dao = mockDao

	req := sub.CreateSubsRequest{
		UserID:   1,
		CourseID: 2,
	}

	resp, err := service.CreateSubs(req)
	assert.NoError(t, err)
	assert.Equal(t, "Subscripcion realizada con exito!", resp.Message)
}

func TestCreateSubs_Error(t *testing.T) {
	mockDao := subs.MockDAO{
		CreateSubsFunc: func(s *sub.Subscription) error {
			return errors.New("DB error")
		},
	}
	service.Dao = mockDao

	req := sub.CreateSubsRequest{
		UserID:   1,
		CourseID: 2,
	}

	_, err := service.CreateSubs(req)
	assert.Error(t, err)
	assert.Equal(t, "DB error", err.Error())
}

func TestGetSubByUserId_Success(t *testing.T) {
	mockDao := subs.MockDAO{
		GetSubByUserIdFunc: func(userId uint) ([]sub.Subscription, error) {
			return []sub.Subscription{
				{ID: 1, UserID: userId, CourseID: 10},
			}, nil
		},
	}
	service.Dao = mockDao

	resp, err := service.GetSubByUserId("1")
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, uint(1), resp[0].ID)
}

func TestGetSubByUserId_InvalidID(t *testing.T) {
	service.Dao = subs.MockDAO{}

	_, err := service.GetSubByUserId("abc")
	assert.Error(t, err)
	assert.Equal(t, "ID invalido", err.Error())
}

func TestGetSubByUserId_ErrorFromDAO(t *testing.T) {
	mockDao := subs.MockDAO{
		GetSubByUserIdFunc: func(userId uint) ([]sub.Subscription, error) {
			return nil, errors.New("DAO error")
		},
	}
	service.Dao = mockDao

	_, err := service.GetSubByUserId("1")
	assert.Error(t, err)
	assert.Equal(t, "DAO error", err.Error())
}

func TestGetSubByCursoId_Success(t *testing.T) {
	mockDao := subs.MockDAO{
		GetSubByCursoIdFunc: func(courseId uint) ([]sub.Subscription, error) {
			return []sub.Subscription{
				{ID: 1, UserID: 1, CourseID: courseId},
				{ID: 2, UserID: 2, CourseID: courseId},
			}, nil
		},
	}
	service.Dao = mockDao

	count, err := service.GetSubByCursoId("1")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestGetSubByCursoId_InvalidID(t *testing.T) {
	service.Dao = subs.MockDAO{}

	_, err := service.GetSubByCursoId("abc")
	assert.Error(t, err)
	assert.Equal(t, "ID invalido", err.Error())
}

func TestGetSubByCursoId_ErrorFromDAO(t *testing.T) {
	mockDao := subs.MockDAO{
		GetSubByCursoIdFunc: func(courseId uint) ([]sub.Subscription, error) {
			return nil, errors.New("DAO error")
		},
	}
	service.Dao = mockDao

	_, err := service.GetSubByCursoId("1")
	assert.Error(t, err)
	assert.Equal(t, "DAO error", err.Error())
}

package subscripcion_service_test

import (
	"subscripciones/models/subs"
	subservice "subscripciones/services/subs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock del DAO
type MockDAO struct {
	mock.Mock
}

func (m *MockDAO) CreateSubs(sub *subs.Subscription) error {
	args := m.Called(sub)
	return args.Error(0)
}

func (m *MockDAO) GetCupos(courseId uint) (subs.GetCuposResp, error) {
	args := m.Called(courseId)
	return args.Get(0).(subs.GetCuposResp), args.Error(1)
}

func (m *MockDAO) GetSubByUserId(userId uint) ([]subs.Subscription, error) {
	args := m.Called(userId)
	return args.Get(0).([]subs.Subscription), args.Error(1)
}

func (m *MockDAO) GetSubByCursoId(courseId uint) ([]subs.Subscription, error) {
	args := m.Called(courseId)
	return args.Get(0).([]subs.Subscription), args.Error(1)
}

func (m *MockDAO) CountSubsByCourse(cursoID int) (int64, error) {
	args := m.Called(cursoID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDAO) DeleteSubByCourseId(courseId uint) error {
	args := m.Called(courseId)
	return args.Error(0)
}

func TestCreateSubs_Success(t *testing.T) {
	mockDAO := new(MockDAO)
	subservice.Dao = mockDAO // Inyectamos el mock

	mockRequest := subs.CreateSubsRequest{
		UserID:   1,
		CourseID: 2,
	}

	mockDAO.On("CreateSubs", mock.Anything).Return(nil)

	response, err := subservice.CreateSubs(mockRequest)

	assert.NoError(t, err)
	assert.Equal(t, "Subscripcion realizada con exito!", response.Message)
	mockDAO.AssertExpectations(t)
}

func TestGetSubByUserId_Success(t *testing.T) {
	mockDAO := new(MockDAO)
	subservice.Dao = mockDAO // Inyectamos el mock

	mockUserID := uint(1)
	mockResponse := []subs.Subscription{
		{ID: 1, UserID: mockUserID, CourseID: 2},
	}

	mockDAO.On("GetSubByUserId", mockUserID).Return(mockResponse, nil)

	response, err := subservice.GetSubByUserId("1")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, mockUserID, response[0].UserID)
	mockDAO.AssertExpectations(t)
}

package subs

import (
	"subscripciones/models/subs"
)

// MockDAO es un mock de SubsDAO
type MockDAO struct {
	CreateSubsFunc        func(*subs.Subscription) error
	GetSubByUserIdFunc    func(uint) ([]subs.Subscription, error)
	GetSubByCursoIdFunc   func(uint) ([]subs.Subscription, error)
}

func (m MockDAO) CreateSubs(sub *subs.Subscription) error {
	return m.CreateSubsFunc(sub)
}

func (m MockDAO) GetCupos(courseId uint) (subs.GetCuposResp, error) {
	return subs.GetCuposResp{}, nil
}

func (m MockDAO) GetSubByUserId(userId uint) ([]subs.Subscription, error) {
	return m.GetSubByUserIdFunc(userId)
}

func (m MockDAO) GetSubByCursoId(courseId uint) ([]subs.Subscription, error) {
	return m.GetSubByCursoIdFunc(courseId)
}

func (m MockDAO) CountSubsByCourse(courseID int) (int64, error) {
	return 0, nil
}

func (m MockDAO) DeleteSubByCourseId(courseId uint) error {
	return nil
}

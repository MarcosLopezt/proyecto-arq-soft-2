package subs

import (
	"subscripciones/models/subs"
)

// SubsDAO define la interfaz para el DAO
type SubsDAO interface {
	CreateSubs(sub *subs.Subscription) error
	GetCupos(courseId uint) (subs.GetCuposResp, error)
	GetSubByUserId(userId uint) ([]subs.Subscription, error)
	GetSubByCursoId(courseId uint) ([]subs.Subscription, error)
	CountSubsByCourse(cursoID int) (int64, error)
	DeleteSubByCourseId(courseId uint) error
}

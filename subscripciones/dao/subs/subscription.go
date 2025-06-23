package subs

import (
	"errors"
	"subscripciones/db"
	subscripciones "subscripciones/models/subs"

	"gorm.io/gorm"
)

// RealDAO es la implementación del DAO real que usa la base de datos
type RealDAO struct{}

func (r RealDAO) CreateSubs(sub *subscripciones.Subscription) error {
	return db.DB.Create(sub).Error
}

func (r RealDAO) GetCupos(courseId uint) (subscripciones.GetCuposResp, error) {
	var cupos subscripciones.GetCuposResp

	if err := db.DB.Where("course_id = ?", courseId).First(&cupos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cupos, errors.New("curso no encontrado")
		}
		return cupos, err
	}

	return cupos, nil
}

func (r RealDAO) GetSubByUserId(userId uint) ([]subscripciones.Subscription, error) {
	var subs []subscripciones.Subscription
	if err := db.DB.Where("user_id = ?", userId).Find(&subs).Error; err != nil {
		return nil, err
	}

	if len(subs) == 0 {
		return nil, errors.New("no subscriptions found for the given user ID")
	}

	return subs, nil
}

func (r RealDAO) GetSubByCursoId(courseId uint) ([]subscripciones.Subscription, error) {
	var subs []subscripciones.Subscription
	if err := db.DB.Where("course_id = ?", courseId).Find(&subs).Error; err != nil {
		return nil, err
	}

	return subs, nil
}

func (r RealDAO) CountSubsByCourse(courseID int) (int64, error) {
	var cant int64
	if err := db.DB.Where("course_id = ?", courseID).Count(&cant).Error; err != nil {
		return 0, err
	}

	return cant, nil
}

func (r RealDAO) DeleteSubByCourseId(courseId uint) error {
	var subs []subscripciones.Subscription

	if err := db.DB.Where("course_id = ?", courseId).Find(&subs).Error; err != nil {
		return err
	}

	if len(subs) > 0 {
		if err := db.DB.Delete(&subs).Error; err != nil {
			return err
		}
	}

	return nil
}

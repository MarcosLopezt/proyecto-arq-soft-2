package subs

import (
	"errors"
	"subscripciones/db"
	subscripciones "subscripciones/models/subs"

	"gorm.io/gorm"
)

func CreateSubs(sub *subscripciones.Subscription) error {
	return db.DB.Create(sub).Error
}

func GetCupos(courseId uint) (subscripciones.GetCuposResp, error) {
	var cupos subscripciones.GetCuposResp
	
	if err := db.DB.Where("course_id = ?", courseId).First(&cupos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cupos, errors.New("curso no encontrado")
		}
		return cupos, err
	}

	return cupos, nil
}

func GetSubByUserId(userId uint) ([]subscripciones.Subscription,error){
	var subs []subscripciones.Subscription
	if err:= db.DB.Where("user_id = ?", userId).Find(&subs).Error; err != nil {
		return nil, err
	}

	if len(subs) == 0 {
        return nil, errors.New("no subscriptions found for the given user ID")
    }

	return subs, nil
}

func GetSubByCursoId(cursoId uint)([]subscripciones.Subscription, error){
	var subs []subscripciones.Subscription
	if err:= db.DB.Where("course_id = ?", cursoId).Find(&subs).Error; err != nil {
		return nil, err
	}
	if len(subs) == 0 {
        return subs, nil  
    }

	return subs, nil
}

func CountSubsByCourse(cursoID int)(int64, error){
	var cant int64
	if err:= db.DB.Where("course_id = ?", cursoID).Count(&cant).Error; err != nil{
		return 0, err
	}

	return cant, nil
}
func DeleteSubByCourseId(courseId uint) error {
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


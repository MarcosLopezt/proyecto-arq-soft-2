package comentarios

import (
	"time"
)

type Comentario struct {
	ID           uint   `json:"ID" gorm:"primaryKey"`
	Valor		 int 	`json:"valor" gorm:"not null"`
	CursoId      uint    `json:"curso_id" gorm:"not null"`
	UserId       uint    `json:"user_id" gorm:"not null"`
	Texto        string       `json:"texto"`
	Fecha        time.Time `gorm:"autoCreateTime"`
}

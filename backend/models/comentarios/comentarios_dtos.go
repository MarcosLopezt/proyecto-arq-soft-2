package comentarios

import (
	"time"
)

type CreateComentRequest struct {
	UserID  uint    `json:"user_id"`
	CursoID uint    `json:"curso_id"`
	Valor	int 	`json:"valor"`
	Texto   string `json:"texto"`
}

type CreateComentResponse struct {
	Message string `json:"message"`
}

type GetComentByCourseReq struct {
	CursoID int `json:"curso_id"`
}

type GetComentByCourseResp struct {
	UserID  uint    `json:"user_id"`
	Valor	int 	`json:"valor"`
	CursoID uint    `json:"curso_id"`
	Texto   string `json:"texto"`
	Fecha  time.Time `json:"fecha"`
}
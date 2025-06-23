package subs

type CreateSubsRequest struct {
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id"`
}

type CreateSubsResponse struct {
	Message string `json:"message"`
}

type GetSubByUserRequest struct {
	UserID uint `json:"user_id"`
}

type GetSubByUserResponse struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id"`
}

type GetCuposReq struct {
	CourseID uint `json:"course_id"`
}

type GetCuposResp struct {
	Cupos int `json:"cupos"`
}
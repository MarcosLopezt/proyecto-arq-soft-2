package files

type UploadFile struct {
	CourseID uint   `json:"curso_id"`
	File     []byte `json:"file"`
}

type UploadFileResponse struct {
	Message string `json:"message"`
}

type GetFileRequest struct {
	CourseID uint `json:"curso_id"`
}

type GetFileResponse struct {
	ID   uint   `json:"ID"`
	File []byte `json:"file"`
}
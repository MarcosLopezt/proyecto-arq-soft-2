package dao

type Curso struct {
	ID          uint   `json:"ID"`
	CourseName  string `json:"course_name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Length      int    `json:"length"`
}

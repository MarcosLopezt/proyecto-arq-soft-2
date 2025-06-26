package cursos

import "time"

type Curso struct {
    ID          uint      `json:"id"`
    CourseName  string    `json:"course_name"`
    Category    string    `json:"category"`
    Length      int       `json:"length"`
    Description string    `json:"description"`
    Cupos       int    `json:"cupos"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CursoNew struct {
    Operation string `json:"operation"`
    CursoID   uint   `json:"course_id"`
}
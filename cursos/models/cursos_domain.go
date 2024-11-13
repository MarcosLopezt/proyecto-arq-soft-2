package cursos

import (
	"time"
)

type Course struct {
    ID          uint   `bson:"_id"`
    CourseName  string    `bson:"course_name"`
    Category    string    `bson:"category"`
    Length      int       `bson:"length"`
    Description string    `bson:"description"`
    CreatedAt   time.Time `bson:"created_at"`  
    UpdatedAt   time.Time `bson:"updated_at"`  
}


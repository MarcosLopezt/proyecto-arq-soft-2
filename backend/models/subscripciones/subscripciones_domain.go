package subscripciones

import (
	"time"
)

type Subscription struct {
    ID           uint      `gorm:"primaryKey"`
    UserID       uint      `gorm:"column:user_id; not null"`
    CourseID     uint      `gorm:"column:course_id; not null"`
    CreationDate time.Time `gorm:"autoCreateTime"`
    LastUpdated  time.Time `gorm:"autoUpdateTime"`
}

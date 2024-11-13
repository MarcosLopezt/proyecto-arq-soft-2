package users

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null;type:varchar(255)"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"not null"`
	CreationDate time.Time `gorm:"autoCreateTime"`
	LastUpdated  time.Time `gorm:"autoUpdateTime"`
}

// func (User) TableName() string {
//     return "user" // Aquí especifica el nombre de la tabla en singular
// }
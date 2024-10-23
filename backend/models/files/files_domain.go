package files

type File struct {
	ID      uint   `gorm:"primaryKey"`
	CursoID uint   `gorm:"not null"`
	File    []byte `gorm:"not null;type:mediumblob"` //nuevos cambios
}
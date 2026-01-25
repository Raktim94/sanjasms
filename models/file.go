package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Filename string `gorm:"not null"`
	Path     string `gorm:"not null"`
	MimeType string
	Size     int64
}

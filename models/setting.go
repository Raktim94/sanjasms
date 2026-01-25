package models

import "gorm.io/gorm"

type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex;not null"`
	Value string
}

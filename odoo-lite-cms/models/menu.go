package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Label    string `gorm:"not null"`
	URL      string `gorm:"not null"`
	Sequence int    `gorm:"default:0"`
	ParentID *uint
	IsActive bool `gorm:"default:true"`
}

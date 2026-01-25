package models

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Slug            string `gorm:"uniqueIndex;not null"`
	Title           string `gorm:"not null"`
	HTMLContent     string // Custom HTML body
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	IsPublished     bool `gorm:"default:false"`
}

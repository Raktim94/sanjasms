package models

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Slug        string `gorm:"uniqueIndex;not null"`
	Title       string `gorm:"not null"`
	HTMLContent string // Custom HTML body

	// SEO Fields
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	OGImage         string // URL to Open Graph Image

	// Custom Code
	CustomHead string // Scripts/Styles for <head>
	CustomBody string // Scripts for end of <body>

	IsPublished bool `gorm:"default:false"`
}

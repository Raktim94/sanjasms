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

	IsPublished bool        `gorm:"default:false"`
	PageType    string      `gorm:"default:'normal'"` // normal, landing
	Blocks      []PageBlock `gorm:"foreignKey:PageID;constraint:OnDelete:CASCADE"`
}

type PageBlock struct {
	gorm.Model
	PageID   uint
	Type     string // text, image, button, embed, html
	Content  string `gorm:"type:text"`
	Sequence int
}

type Embed struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Code     string `gorm:"not null"` // HTML/JS/IFrame code
	IsActive bool   `gorm:"default:true"`
}

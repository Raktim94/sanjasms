package models

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbPath := os.Getenv("DB_Path")
	if dbPath == "" {
		dbPath = "data/cms.db"
	}

	// Ensure directory exists
	// os.MkdirAll("data", 0755) // Docker setup handles /data volume, but good for local dev

	var database *gorm.DB
	var err error

	// Retry logic for Docker startup racing
	for i := 0; i < 10; i++ {
		database, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/10): %v. Retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}

	err = database.AutoMigrate(&User{}, &Page{}, &Menu{}, &Setting{}, &File{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database

	// Initialize Default Data
	initializeDefaults(database)
}

func initializeDefaults(db *gorm.DB) {
	// 1. Default User
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("hi@raktim"), bcrypt.DefaultCost)
		user := User{
			Email:        "hi@RAKTIMranjit.in",
			PasswordHash: string(hash),
			IsAdmin:      true,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create default user: %v", err)
		} else {
			log.Println("Default user 'hi@RAKTIMranjit.in' created.")
		}
	}

	// 2. Default Settings
	defaults := map[string]string{
		"site_title":       "sanjanacms",
		"footer_text":      "Odoo Lite © 2026 Odoo Lite CMS",
		"primary_color":    "#007AFF", // Apple Blue
		"meta_description": "A lightweight CMS",
	}

	for key, value := range defaults {
		var setting Setting
		if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
			db.Create(&Setting{Key: key, Value: value})
		}
	}
}

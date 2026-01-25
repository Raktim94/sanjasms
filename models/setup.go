package models

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
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
}

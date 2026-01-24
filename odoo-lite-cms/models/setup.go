package models

import (
	"log"
	"os"

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

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = database.AutoMigrate(&User{}, &Page{}, &Menu{}, &Setting{}, &File{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database
}

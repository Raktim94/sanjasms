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

	log.Printf("Initializing database connection at: %s", dbPath)

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
		log.Printf("CRITICAL: Failed to connect to database after retries: %v", err)
		// log.Fatal would exit, but we want to see if we can keep container alive to read logs if needed,
		// but app needs DB. So Fatal is appropriate, but let's make sure it's flushed.
		log.Fatal(err)
	}

	err = database.AutoMigrate(&User{}, &Page{}, &Menu{}, &Setting{}, &File{}, &Embed{}, &PageBlock{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database
	log.Println("Database connection established and migrated.")

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
		"footer_text":      "sanjanacms © 2026 sanjanacms",
		"primary_color":    "#007AFF", // Apple Blue
		"meta_description": "A lightweight CMS",
		"robots_txt":       "User-agent: *\nAllow: /",
	}

	for key, value := range defaults {
		var setting Setting
		if err := db.Where("key = ?", key).First(&setting).Error; err != nil {
			db.Create(&Setting{Key: key, Value: value})
		}
	}

	// 3. Demo Pages
	demoPages := []Page{
		{
			Slug:            "home",
			Title:           "Home",
			MetaTitle:       "Home - sanjanacms",
			MetaDescription: "Welcome to sanjanacms demo site.",
			HTMLContent: `
				<section class="hero" style="text-align: center; padding: 100px 20px;">
					<h1 style="font-size: 3rem; margin-bottom: 20px;">Welcome to sanjanacms</h1>
					<p style="font-size: 1.2rem; color: #666; max-width: 600px; margin: 0 auto 30px;">
						A simple, powerful, and lightweight content management system.
						This is a demo page you can edit instantly.
					</p>
					<a href="/about" class="btn-cta">Learn More</a>
				</section>
				<section class="features" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 40px; padding: 60px 20px; max-width: 1200px; margin: 0 auto;">
					<div class="feature">
						<h3>🚀 Fast Performance</h3>
						<p>Pages load instantly with optimized rendering.</p>
					</div>
					<div class="feature">
						<h3>✏️ Easy Editing</h3>
						<p>Manage your content with a simple dashboard.</p>
					</div>
					<div class="feature">
						<h3>🔐 Secure & Safe</h3>
						<p>Built with security in mind for peace of mind.</p>
					</div>
				</section>
			`,
			IsPublished: true,
			IsDraft:     false,
		},
		{
			Slug:            "about",
			Title:           "About Us",
			MetaTitle:       "About - sanjanacms",
			MetaDescription: "Learn more about our mission.",
			HTMLContent: `
				<div style="max-width: 800px; margin: 0 auto; padding: 60px 20px;">
					<h1>About Us</h1>
					<p>We are dedicated to building the best lightweight CMS experience.</p>
					<p>This platform allows you to manage content without the bloat.</p>
				</div>
			`,
			IsPublished: true,
			IsDraft:     false,
		},
		{
			Slug:            "contact",
			Title:           "Contact",
			MetaTitle:       "Contact Us - sanjanacms",
			MetaDescription: "Get in touch with us.",
			HTMLContent: `
				<div style="max-width: 800px; margin: 0 auto; padding: 60px 20px;">
					<h1>Contact Us</h1>
					<p>Email us at: <a href="mailto:contact@example.com">contact@example.com</a></p>
				</div>
			`,
			IsPublished: true,
			IsDraft:     false,
		},
	}

	for _, p := range demoPages {
		var page Page
		if err := db.Where("slug = ?", p.Slug).First(&page).Error; err != nil {
			db.Create(&p)
			log.Printf("Created demo page: %s", p.Slug)
		} else {
			// Ensure essential demo pages are published if they exist (fixes partial seed issues)
			if !page.IsPublished || page.IsDraft {
				page.IsPublished = true
				page.IsDraft = false
				db.Save(&page)
				log.Printf("Corrected demo page to published and non-draft: %s", p.Slug)
			}
		}
	}

	// 4. Demo Menus
	var menuCount int64
	db.Model(&Menu{}).Count(&menuCount)
	if menuCount == 0 {
		menus := []Menu{
			{Label: "Home", URL: "/", Sequence: 1, IsActive: true},
			{Label: "About", URL: "/about", Sequence: 2, IsActive: true},
			{Label: "Contact", URL: "/contact", Sequence: 3, IsActive: true},
		}
		for _, m := range menus {
			db.Create(&m)
		}
		log.Println("Created default menus.")
	}
}

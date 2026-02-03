package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"odoo-lite-cms/handlers"
	"odoo-lite-cms/middleware"
	"odoo-lite-cms/models"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Custom panic handler for debugging
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					log.Printf("PANIC RECOVERED: %v\nStack: %s", err, debug.Stack())
					c.Error(err)
				}
			}()
			return next(c)
		}
	})
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Session Middleware (Secure Cookie)
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "super-secret-key-change-it-in-production"
	}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))

	// CSRF Protection
	e.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		TokenLookup: "form:_csrf,header:X-CSRF-Token", // Look for CSRF token in hidden form input or header
		CookiePath:  "/",
	}))

	// Rate Limiting
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20))) // 20 requests per second

	// CORS (Optional but good for API consistency)
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Template Registry
	t := &Template{
		templates: make(map[string]*template.Template),
	}

	// Helper functions
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// 1. Identify Layouts and Partials
	log.Println("Loading templates...")
	commonTemplates := []string{
		"templates/admin/layout.html",
	}

	// 2. Identify Page Templates
	loadTemplates := func() {
		// Admin Pages
		log.Println("Globbing templates/admin/...")
		adminFiles, err := filepath.Glob("templates/admin/*.html")
		if err != nil {
			log.Printf("ERROR globbing admin templates: %v", err)
		}

		for _, file := range adminFiles {
			// Skip layout files if we don't want them as standalone pages
			if strings.Contains(file, "layout.html") {
				continue
			}

			// Create a new template set for this page
			fileName := filepath.Base(file)
			// Parse layouts first
			tmpl := template.New(fileName).Funcs(funcMap)
			// Parse common layouts
			if len(commonTemplates) > 0 {
				_, err = tmpl.ParseFiles(commonTemplates...)
				if err != nil {
					log.Printf("Error parsing layouts for %s: %v", fileName, err)
					continue
				}
			}
			// Parse the page itself
			_, err = tmpl.ParseFiles(file)
			if err != nil {
				log.Printf("Error parsing file %s: %v", fileName, err)
				continue
			}

			t.templates[fileName] = tmpl
			log.Printf("Registered template: %s", fileName)
		}

		// Public Pages
		log.Println("Globbing templates/public/...")
		publicLayout := "templates/public/layout.html"
		publicFiles, err := filepath.Glob("templates/public/*.html")
		if err != nil {
			log.Printf("ERROR globbing public templates: %v", err)
		} else {
			for _, file := range publicFiles {
				// We need to parse layout.html as itself to make "public_layout" available
				// But we also need to parse it alongside page.html so page.html can "see" it.

				if strings.Contains(file, "layout.html") {
					continue
				}

				fileName := filepath.Base(file)
				tmpl := template.New(fileName).Funcs(funcMap)

				// Parse layout FIRST, then the file.
				// The order matters. layout.html defines "public_layout".
				filesToParse := []string{publicLayout, file}

				// Verify layout exists
				if _, err := os.Stat(publicLayout); err != nil {
					log.Fatalf("FATAL ERROR: Public layout layout.html missing at %s", publicLayout)
				}

				if _, err := tmpl.ParseFiles(filesToParse...); err != nil {
					// CRITICAL: Fail fast if templates are broken
					log.Fatalf("FATAL ERROR: Could not parse public template %s: %v", fileName, err)
				}

				t.templates[fileName] = tmpl
				log.Printf("SUCCESS: Registered public template: %s (Linked to public_layout)", fileName)
			}

			// Also register the 404 page explicitly if it wasn't caught above
			if _, ok := t.templates["404.html"]; !ok {
				tmpl := template.New("404.html").Funcs(funcMap)
				if _, err := tmpl.ParseFiles(publicLayout, "templates/public/404.html"); err == nil {
					t.templates["404.html"] = tmpl
					log.Printf("Registered 404 template explicitly")
				}
			}
		}

		// Root templates
		log.Println("Globbing templates/ (root)...")
		rootFiles, _ := filepath.Glob("templates/*.html")
		for _, file := range rootFiles {
			fileName := filepath.Base(file)
			tmpl := template.New(fileName).Funcs(funcMap)
			tmpl.ParseFiles(file)
			t.templates[fileName] = tmpl
			log.Printf("Registered template: %s", fileName)
		}
	}

	loadTemplates()

	e.Renderer = t

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	log.Println("About to connect to database...")
	models.ConnectDatabase()

	// Static Setup
	e.Static("/static", "static")
	e.Static("/uploads", "static/uploads")

	// Routes
	// Auth Routes
	e.GET("/admin/login", handlers.LoginPage)
	e.POST("/admin/login", handlers.LoginService)
	e.GET("/admin/logout", handlers.Logout)
	e.GET("/admin/setup", handlers.SetupPage)
	e.POST("/admin/setup", handlers.SetupService)

	// Admin Routes (Protected)
	adminGroup := e.Group("/admin")
	e.GET("/admin", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
	})
	adminGroup.Use(middleware.RequireAuth)
	adminGroup.GET("/dashboard", handlers.AdminDashboard)

	// Page Routes
	adminGroup.GET("/pages", handlers.PageList)
	adminGroup.GET("/pages/new", handlers.PageNew)
	adminGroup.POST("/pages", handlers.PageCreate)
	adminGroup.GET("/pages/:id/edit", handlers.PageEdit)
	adminGroup.POST("/pages/:id", handlers.PageUpdate)
	adminGroup.POST("/pages/:id/delete", handlers.PageDelete)

	// Menu Routes
	adminGroup.GET("/menus", handlers.MenuList)
	adminGroup.GET("/menus/new", handlers.MenuNew)
	adminGroup.POST("/menus", handlers.MenuCreate)
	adminGroup.GET("/menus/:id/edit", handlers.MenuEdit)
	adminGroup.POST("/menus/:id", handlers.MenuUpdate)
	adminGroup.POST("/menus/:id/delete", handlers.MenuDelete)
	adminGroup.POST("/menus/reorder", handlers.MenuReorder)

	// Media Routes
	adminGroup.GET("/media", handlers.MediaList)
	adminGroup.POST("/media/upload", handlers.MediaUpload)
	adminGroup.POST("/media/:id/delete", handlers.MediaDelete)

	// Embed Routes
	adminGroup.GET("/embeds", handlers.EmbedList)
	adminGroup.GET("/embeds/new", handlers.EmbedNew)
	adminGroup.POST("/embeds", handlers.EmbedCreate)
	adminGroup.GET("/embeds/:id/edit", handlers.EmbedEdit)
	adminGroup.POST("/embeds/:id", handlers.EmbedUpdate)
	adminGroup.POST("/embeds/:id/delete", handlers.EmbedDelete)

	// User Routes
	adminGroup.GET("/users", handlers.UserList)
	adminGroup.GET("/users/new", handlers.UserNew)
	adminGroup.POST("/users", handlers.UserCreate)
	adminGroup.POST("/users/:id/delete", handlers.UserDelete)

	// Settings Routes
	adminGroup.GET("/settings", handlers.SettingsPage)
	adminGroup.POST("/settings", handlers.SettingsUpdate)

	// Public Sitemap
	e.GET("/sitemap.xml", handlers.Sitemap)
	e.GET("/robots.txt", func(c echo.Context) error {
		settings := handlers.GetSettingsMap()
		return c.String(http.StatusOK, settings["robots_txt"])
	})

	// Catch-all for public pages
	e.GET("/*", handlers.PublicPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server with Graceful Shutdown
	go func() {
		log.Printf("Starting sanjanacms on port %s...", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

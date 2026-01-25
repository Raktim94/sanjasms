package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"odoo-lite-cms/handlers"
	"odoo-lite-cms/middleware"
	"odoo-lite-cms/models"

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
	// We execute the "layout" template, which should invoke the specific page content
	// Actually, wait. The structure is usually: layout defines "content" block.
	// If we use standard Go templates, we need to execute the *layout* template name if defined,
	// or the file name if it's the entry point.
	// In our Setup:
	// layout.html defines "admin_layout".
	// dashboard.html defines "dashboard.html" and "content".
	// dashboard.html SAYS {{template "admin_layout" .}}
	// So we should execute "dashboard.html" (the name of the file/block).
	return tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Generic middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())

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
	// We assume templates/admin/layout.html is the main layout for admin.
	// We might have others.
	// Strategic approach:
	// - Parse all layouts/partials first.
	// - For EACH page template, clone the set or create a new set including layouts + that page.

	// Let's define common layouts.
	commonTemplates := []string{
		"templates/admin/layout.html",
		// Add other shared components here if any
	}

	// 2. Identify Page Templates
	// specific pages in templates/admin/
	// We want to exclude layout.html from being treated as a "page" textually, though it's in the dir.

	// Helper to load templates
	loadTemplates := func() {
		// Admin Pages
		adminFiles, err := filepath.Glob("templates/admin/*.html")
		if err != nil {
			log.Fatal(err)
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

		// Public Pages? (Similar logic if they use layouts)
		// For now, let's look at public/layout.html and public/*.html
		publicCommon := []string{"templates/public/layout.html"}
		publicFiles, err := filepath.Glob("templates/public/*.html")
		if err != nil {
			log.Println("No public templates found or error")
		} else {
			for _, file := range publicFiles {
				if strings.Contains(file, "layout.html") {
					continue
				}
				fileName := filepath.Base(file)
				tmpl := template.New(fileName).Funcs(funcMap)
				if len(publicCommon) > 0 {
					// Check existence because public layout might not exist yet
					if _, err := os.Stat(publicCommon[0]); err == nil {
						tmpl.ParseFiles(publicCommon...)
					}
				}
				tmpl.ParseFiles(file)
				t.templates[fileName] = tmpl
				log.Printf("Registered template: %s", fileName)
			}
		}

		// Root templates (e.g. login.html if it's in root)
		// Assuming templates/login.html exists and might not use the admin layout?
		// Or if it does. Let's see.
		// If templates/*.html exists
		rootFiles, _ := filepath.Glob("templates/*.html")
		for _, file := range rootFiles {
			fileName := filepath.Base(file)
			// Maybe these stand alone?
			tmpl := template.New(fileName).Funcs(funcMap)
			tmpl.ParseFiles(file)
			t.templates[fileName] = tmpl
			log.Printf("Registered template: %s", fileName)
		}
	}

	loadTemplates()

	e.Renderer = t

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

	// Media Routes
	adminGroup.GET("/media", handlers.MediaList)
	adminGroup.POST("/media/upload", handlers.MediaUpload)
	adminGroup.POST("/media/:id/delete", handlers.MediaDelete)

	// Settings Routes
	adminGroup.GET("/settings", handlers.SettingsPage)
	adminGroup.POST("/settings", handlers.SettingsUpdate)

	// Catch-all for public pages
	e.GET("/*", handlers.PublicPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting sanjanacms on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}

package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"odoo-lite-cms/handlers"
	"odoo-lite-cms/middleware"
	"odoo-lite-cms/models"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Register Template Renderer
	// We recursively parse templates. For simplicity now, just flattened logic or explicit patterns
	// Note: ParseGlob doesn't recurse directory. We need multiple globs or Walk.
	// For now, let's just parse specific patterns we know we created.
	t := &Template{
		templates: template.New("main").Funcs(template.FuncMap{
			"safe": func(s string) template.HTML {
				return template.HTML(s)
			},
		}),
	}
	// We will parse admin/ later when we create them.
	// Ideally we should parse both.
	// t.templates.ParseGlob("templates/admin/*.html")
	// t.templates.ParseGlob("templates/public/*.html")

	// Improving template parsing to include subdirectories
	func() {
		// Parse root templates
		// We need to parse all files into the same set to allow cross-referencing if needed,
		// or just distinct logic.
		// Let's assume unique names for now.
		patterns := []string{"templates/*.html", "templates/admin/*.html", "templates/public/*.html"}
		for _, pattern := range patterns {
			// loose check if files exist to avoid panic on empty folders
			matches, _ := filepath.Glob(pattern)
			if len(matches) > 0 {
				template.Must(t.templates.ParseGlob(pattern))
			}
		}
	}()

	e.Renderer = t

	models.ConnectDatabase()

	// Static Setup
	e.Static("/static", "static")
	e.Static("/uploads", "static/uploads")

	// Public Routes (Catch-all for pages)
	// We register this last? No, Echo matches in order of specificity usually, or declaration order.
	// Static and Admin are specific.
	// "/*" will match everything not matched.
	// However, explicit paths take precedence. "/*" might shadow specific if defined earlier.
	// Let's rely on echo's router logic.
	// e.GET("/", ... ) -> defined earlier was simple string.

	// We will use a catch-all handler for slug-based pages
	e.GET("/*", handlers.PublicPage)

	// Auth Routes
	e.GET("/admin/login", handlers.LoginPage)
	e.POST("/admin/login", handlers.LoginService)
	e.GET("/admin/logout", handlers.Logout)
	e.GET("/admin/setup", handlers.SetupPage)
	e.POST("/admin/setup", handlers.SetupService)

	// Admin Routes (Protected)
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.RequireAuth)
	adminGroup.GET("/dashboard", handlers.AdminDashboard) // Need to implement this

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}

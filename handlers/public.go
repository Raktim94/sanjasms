package handlers

import (
	"html/template"
	"log"
	"net/http"
	"odoo-lite-cms/models"
	"strings"

	"github.com/labstack/echo/v4"
)

func PublicPage(c echo.Context) error {
	slug := c.Param("*")
	// If slug is empty or "/", lookup "home" or first page.
	if slug == "" || slug == "/" {
		slug = "home"
	}
	// Clean leading slash
	slug = strings.TrimPrefix(slug, "/")

	log.Printf("PublicPage requested: slug='%s'", slug)

	// Safe retrieval of shared data
	menus := GetPublicMenus()
	if menus == nil {
		menus = []models.Menu{}
	}
	settings := GetSettingsMap()
	if settings == nil {
		settings = map[string]string{
			"site_title":  "sanjanacms",
			"footer_text": "Default Footer",
		}
	}

	var page models.Page
	// Find published page by slug
	if err := models.DB.Where("slug = ? AND is_published = ?", slug, true).First(&page).Error; err != nil {
		log.Printf("Page not found in DB: %s (error: %v)", slug, err)

		// Render 404 with safe context
		// IMPORTANT: Pass a nil Page to force template fallback logic
		return c.Render(http.StatusNotFound, "404.html", map[string]interface{}{
			"Title":    "Page Not Found",
			"Settings": settings,
			"Menus":    menus,
			"Page":     nil,
		})
	}

	log.Printf("Rendering page: %s with template 'page.html'", page.Title)

	err := c.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Title":           page.MetaTitle, // Fallback to Title if empty handled in template
		"MetaDescription": page.MetaDescription,
		"MetaKeywords":    page.MetaKeywords,
		"Content":         template.HTML(page.HTMLContent), // Trusting admin content as safe
		"Page":            page,
		"Menus":           menus,
		"Settings":        settings,
	})

	if err != nil {
		log.Printf("ERROR rendering template 'page.html': %v", err)
		return err
	}
	return nil
}

func GetPublicMenus() []models.Menu {
	var menus []models.Menu
	if models.DB == nil {
		return menus
	}
	models.DB.Where("is_active = ?", true).Order("sequence asc").Find(&menus)
	return menus
}

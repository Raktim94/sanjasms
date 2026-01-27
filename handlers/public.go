package handlers

import (
	"html/template"
	"log"
	"net/http"
	"odoo-lite-cms/models"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	if err := models.DB.Preload("Blocks", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Where("slug = ? AND is_published = ?", slug, true).First(&page).Error; err != nil {
		log.Printf("Page not found in DB: %s (error: %v)", slug, err)

		// Render 404 with safe context
		var embeds []models.Embed
		models.DB.Where("is_active = ?", true).Find(&embeds)
		return c.Render(http.StatusNotFound, "404.html", map[string]interface{}{
			"Title":    "Page Not Found",
			"Settings": settings,
			"Menus":    menus,
			"Page":     nil,
			"Embeds":   embeds,
		})
	}

	var embeds []models.Embed
	models.DB.Where("is_active = ?", true).Find(&embeds)

	log.Printf("Rendering page: %s with template 'page.html'", page.Title)

	err := c.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Title":           page.MetaTitle, // Fallback to Title if empty handled in template
		"MetaDescription": page.MetaDescription,
		"MetaKeywords":    page.MetaKeywords,
		"Content":         template.HTML(page.HTMLContent), // Trusting admin content as safe
		"Page":            page,
		"Menus":           menus,
		"Settings":        settings,
		"Embeds":          embeds,
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

func Sitemap(c echo.Context) error {
	var pages []models.Page
	if err := models.DB.Where("is_published = ?", true).Find(&pages).Error; err != nil {
		return c.XML(http.StatusInternalServerError, nil)
	}

	urlSet := struct {
		XMLName string `xml:"urlset"`
		Xmlns   string `xml:"xmlns,attr"`
		URLs    []struct {
			Loc        string `xml:"loc"`
			LastMod    string `xml:"lastmod"`
			ChangeFreq string `xml:"changefreq"`
			Priority   string `xml:"priority"`
		} `xml:"url"`
	}{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	baseURL := "http://" + c.Request().Host // Dynamic base URL

	// Add Home
	urlSet.URLs = append(urlSet.URLs, struct {
		Loc        string `xml:"loc"`
		LastMod    string `xml:"lastmod"`
		ChangeFreq string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	}{
		Loc:        baseURL + "/",
		ChangeFreq: "daily",
		Priority:   "1.0",
	})

	for _, p := range pages {
		if p.Slug == "home" {
			continue
		}
		urlSet.URLs = append(urlSet.URLs, struct {
			Loc        string `xml:"loc"`
			LastMod    string `xml:"lastmod"`
			ChangeFreq string `xml:"changefreq"`
			Priority   string `xml:"priority"`
		}{
			Loc:        baseURL + "/" + p.Slug,
			ChangeFreq: "weekly",
			Priority:   "0.8",
		})
	}

	return c.XML(http.StatusOK, urlSet)
}

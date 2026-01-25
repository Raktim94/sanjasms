package handlers

import (
	"net/http"
	"odoo-lite-cms/models"

	"github.com/labstack/echo/v4"
)

func PageList(c echo.Context) error {
	var pages []models.Page
	models.DB.Find(&pages)
	return c.Render(http.StatusOK, "pages_list.html", map[string]interface{}{
		"Title":  "Pages",
		"Active": "pages",
		"Pages":  pages,
	})
}

func PageNew(c echo.Context) error {
	return c.Render(http.StatusOK, "pages_form.html", map[string]interface{}{
		"Title":  "New Page",
		"Active": "pages",
		"Page":   models.Page{},
		"IsNew":  true,
	})
}

func PageCreate(c echo.Context) error {
	page := models.Page{
		Title:           c.FormValue("title"),
		Slug:            c.FormValue("slug"),
		HTMLContent:     c.FormValue("html_content"),
		MetaTitle:       c.FormValue("meta_title"),
		MetaDescription: c.FormValue("meta_description"),
		MetaKeywords:    c.FormValue("meta_keywords"),
		OGImage:         c.FormValue("og_image"),
		CustomHead:      c.FormValue("custom_head"),
		CustomBody:      c.FormValue("custom_body"),
		IsPublished:     c.FormValue("is_published") == "on",
	}

	if page.Slug == "" {
		page.Slug = GenerateSlug(page.Title)
	}

	if err := models.DB.Create(&page).Error; err != nil {
		return c.Render(http.StatusOK, "pages_form.html", map[string]interface{}{
			"Title":  "New Page",
			"Active": "pages",
			"Page":   page,
			"Error":  "Could not save page. ensure slug is unique.",
			"IsNew":  true,
		})
	}
	return c.Redirect(http.StatusFound, "/admin/pages")
}

func PageEdit(c echo.Context) error {
	id := c.Param("id")
	var page models.Page
	if err := models.DB.First(&page, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Page not found")
	}
	return c.Render(http.StatusOK, "pages_form.html", map[string]interface{}{
		"Title":  "Edit Page",
		"Active": "pages",
		"Page":   page,
		"IsNew":  false,
	})
}

func PageUpdate(c echo.Context) error {
	id := c.Param("id")
	var page models.Page
	if err := models.DB.First(&page, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Page not found")
	}

	page.Title = c.FormValue("title")
	page.Slug = c.FormValue("slug")
	page.HTMLContent = c.FormValue("html_content")
	page.MetaTitle = c.FormValue("meta_title")
	page.MetaDescription = c.FormValue("meta_description")
	page.MetaKeywords = c.FormValue("meta_keywords")
	page.OGImage = c.FormValue("og_image")
	page.CustomHead = c.FormValue("custom_head")
	page.CustomBody = c.FormValue("custom_body")
	page.IsPublished = c.FormValue("is_published") == "on"

	if err := models.DB.Save(&page).Error; err != nil {
		return c.Render(http.StatusOK, "pages_form.html", map[string]interface{}{
			"Title":  "Edit Page",
			"Active": "pages",
			"Page":   page,
			"Error":  "Could not update page.",
			"IsNew":  false,
		})
	}
	return c.Redirect(http.StatusFound, "/admin/pages")
}

func PageDelete(c echo.Context) error {
	id := c.Param("id")
	models.DB.Delete(&models.Page{}, id)
	return c.Redirect(http.StatusFound, "/admin/pages")
}

// Simple slug generator helper (placeholder)
// Simple slug generator helper
func GenerateSlug(title string) string {
	// Simple replacement: spaces -> dashes, lowercase
	// Production apps should use a library like goslugify
	slug := ""
	for _, c := range title {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			slug += string(c)
		} else if c >= 'A' && c <= 'Z' {
			slug += string(c + 32) // to lower
		} else if c == ' ' || c == '-' {
			slug += "-"
		}
	}
	// TODO: Verify Uniqueness in DB loop
	return slug
}

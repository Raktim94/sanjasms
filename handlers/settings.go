package handlers

import (
	"log"
	"net/http"
	"odoo-lite-cms/models"

	"github.com/labstack/echo/v4"
)

// Helper to get settings as map
func GetSettingsMap() map[string]string {
	settingsMap := make(map[string]string)

	if models.DB == nil {
		return settingsMap
	}

	var settings []models.Setting
	if err := models.DB.Find(&settings).Error; err != nil {
		log.Printf("Error fetching settings: %v", err)
		return settingsMap
	}

	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}
	return settingsMap
}

func SettingsPage(c echo.Context) error {
	settings := GetSettingsMap()
	return c.Render(http.StatusOK, "settings.html", map[string]interface{}{
		"Title":    "Site Settings",
		"Active":   "settings",
		"Settings": settings,
		"csrf":     c.Get("csrf"),
	})
}

func SettingsUpdate(c echo.Context) error {
	keys := []string{"site_title", "logo_url", "favicon_url", "primary_color", "footer_text", "meta_description", "robots_txt", "allow_public_signup"}

	for _, key := range keys {
		value := c.FormValue(key)
		// Special handling for checkbox
		if key == "allow_public_signup" && value == "" {
			value = "false"
		}

		var setting models.Setting
		if err := models.DB.Where("key = ?", key).First(&setting).Error; err != nil {
			setting = models.Setting{Key: key, Value: value}
			models.DB.Create(&setting)
		} else {
			setting.Value = value
			models.DB.Save(&setting)
		}
	}

	return c.Redirect(http.StatusFound, "/admin/settings")
}

package handlers

import (
	"net/http"
	"odoo-lite-cms/models"

	"github.com/labstack/echo/v4"
)

func EmbedList(c echo.Context) error {
	var embeds []models.Embed
	models.DB.Find(&embeds)
	return c.Render(http.StatusOK, "embeds_list.html", map[string]interface{}{
		"Title":  "Embed Manager",
		"Embeds": embeds,
		"Active": "embeds",
	})
}

func EmbedNew(c echo.Context) error {
	return c.Render(http.StatusOK, "embeds_form.html", map[string]interface{}{
		"Title":  "New Embed",
		"Active": "embeds",
	})
}

func EmbedCreate(c echo.Context) error {
	embed := new(models.Embed)
	if err := c.Bind(embed); err != nil {
		return err
	}
	isActive := c.FormValue("is_active") == "on"
	embed.IsActive = isActive

	if err := models.DB.Create(&embed).Error; err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/admin/embeds")
}

func EmbedEdit(c echo.Context) error {
	id := c.Param("id")
	var embed models.Embed
	if err := models.DB.First(&embed, id).Error; err != nil {
		return err
	}
	return c.Render(http.StatusOK, "embeds_form.html", map[string]interface{}{
		"Title":  "Edit Embed",
		"Embed":  embed,
		"Active": "embeds",
	})
}

func EmbedUpdate(c echo.Context) error {
	id := c.Param("id")
	var embed models.Embed
	if err := models.DB.First(&embed, id).Error; err != nil {
		return err
	}
	if err := c.Bind(&embed); err != nil {
		return err
	}
	isActive := c.FormValue("is_active") == "on"
	embed.IsActive = isActive

	models.DB.Save(&embed)
	return c.Redirect(http.StatusSeeOther, "/admin/embeds")
}

func EmbedDelete(c echo.Context) error {
	id := c.Param("id")
	models.DB.Delete(&models.Embed{}, id)
	return c.Redirect(http.StatusSeeOther, "/admin/embeds")
}

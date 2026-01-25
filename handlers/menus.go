package handlers

import (
	"net/http"
	"odoo-lite-cms/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func MenuList(c echo.Context) error {
	var menus []models.Menu
	models.DB.Order("sequence asc").Find(&menus)
	return c.Render(http.StatusOK, "menus_list.html", map[string]interface{}{
		"Title":  "Menus",
		"Active": "menus",
		"Menus":  menus,
	})
}

func MenuNew(c echo.Context) error {
	return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
		"Title":  "New Menu Item",
		"Active": "menus",
		"Menu":   models.Menu{IsActive: true, Sequence: 10},
		"IsNew":  true,
	})
}

func MenuCreate(c echo.Context) error {
	seq, _ := strconv.Atoi(c.FormValue("sequence"))
	menu := models.Menu{
		Label:    c.FormValue("label"),
		URL:      c.FormValue("url"),
		Sequence: seq,
		IsActive: c.FormValue("is_active") == "on",
	}

	if err := models.DB.Create(&menu).Error; err != nil {
		return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
			"Title":  "New Menu Item",
			"Active": "menus",
			"Menu":   menu,
			"Error":  "Could not save menu item.",
			"IsNew":  true,
		})
	}
	return c.Redirect(http.StatusFound, "/admin/menus")
}

func MenuEdit(c echo.Context) error {
	id := c.Param("id")
	var menu models.Menu
	if err := models.DB.First(&menu, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Menu item not found")
	}
	return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
		"Title":  "Edit Menu Item",
		"Active": "menus",
		"Menu":   menu,
		"IsNew":  false,
	})
}

func MenuUpdate(c echo.Context) error {
	id := c.Param("id")
	var menu models.Menu
	if err := models.DB.First(&menu, id).Error; err != nil {
		return c.String(http.StatusNotFound, "Menu item not found")
	}

	menu.Label = c.FormValue("label")
	menu.URL = c.FormValue("url")
	seq, _ := strconv.Atoi(c.FormValue("sequence"))
	menu.Sequence = seq
	menu.IsActive = c.FormValue("is_active") == "on"

	if err := models.DB.Save(&menu).Error; err != nil {
		return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
			"Title":  "Edit Menu Item",
			"Active": "menus",
			"Menu":   menu,
			"Error":  "Could not update menu item.",
			"IsNew":  false,
		})
	}
	return c.Redirect(http.StatusFound, "/admin/menus")
}

func MenuDelete(c echo.Context) error {
	id := c.Param("id")
	models.DB.Delete(&models.Menu{}, id)
	return c.Redirect(http.StatusFound, "/admin/menus")
}

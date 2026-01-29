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
	var allMenus []models.Menu
	models.DB.Order("label asc").Find(&allMenus)

	return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
		"Title":    "New Menu Item",
		"Active":   "menus",
		"Menu":     models.Menu{IsActive: true, Sequence: 10},
		"IsNew":    true,
		"AllMenus": allMenus,
		"csrf":     c.Get("csrf"),
	})
}

func MenuCreate(c echo.Context) error {
	seq, _ := strconv.Atoi(c.FormValue("sequence"))

	var parentID *uint
	if pid := c.FormValue("parent_id"); pid != "" {
		id, _ := strconv.Atoi(pid)
		u := uint(id)
		parentID = &u
	}

	menu := models.Menu{
		Label:    c.FormValue("label"),
		URL:      c.FormValue("url"),
		Sequence: seq,
		IsActive: c.FormValue("is_active") == "on",
		ParentID: parentID,
	}

	if err := models.DB.Create(&menu).Error; err != nil {
		return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
			"Title":  "New Menu Item",
			"Active": "menus",
			"Menu":   menu,
			"Error":  "Could not save menu item.",
			"IsNew":  true,
			"csrf":   c.Get("csrf"),
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

	var allMenus []models.Menu
	models.DB.Order("label asc").Find(&allMenus)

	return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
		"Title":    "Edit Menu Item",
		"Active":   "menus",
		"Menu":     menu,
		"IsNew":    false,
		"AllMenus": allMenus,
		"csrf":     c.Get("csrf"),
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

	var parentID *uint
	if pid := c.FormValue("parent_id"); pid != "" {
		id, _ := strconv.Atoi(pid)
		u := uint(id)
		parentID = &u
	}
	menu.ParentID = parentID

	if err := models.DB.Save(&menu).Error; err != nil {
		return c.Render(http.StatusOK, "menus_form.html", map[string]interface{}{
			"Title":  "Edit Menu Item",
			"Active": "menus",
			"Menu":   menu,
			"Error":  "Could not update menu item.",
			"IsNew":  false,
			"csrf":   c.Get("csrf"),
		})
	}
	return c.Redirect(http.StatusFound, "/admin/menus")
}

func MenuDelete(c echo.Context) error {
	id := c.Param("id")
	models.DB.Delete(&models.Menu{}, id)
	return c.Redirect(http.StatusFound, "/admin/menus")
}

func MenuReorder(c echo.Context) error {
	var data []struct {
		ID       uint  `json:"id"`
		Sequence int   `json:"sequence"`
		ParentID *uint `json:"parent_id"`
	}

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data"})
	}

	for _, item := range data {
		models.DB.Model(&models.Menu{}).Where("id = ?", item.ID).Updates(map[string]interface{}{
			"sequence":  item.Sequence,
			"parent_id": item.ParentID,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

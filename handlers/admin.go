package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
		"Title":  "Dashboard",
		"Active": "dashboard",
		"csrf":   c.Get("csrf"),
	})
}

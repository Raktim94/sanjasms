package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Simple cookie-based auth for now
		_, err := c.Cookie("admin_session")
		if err != nil {
			return c.Redirect(http.StatusFound, "/admin/login")
		}
		return next(c)
	}
}

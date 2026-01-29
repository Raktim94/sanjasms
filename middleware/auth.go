package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("admin_session", c)
		if err != nil {
			log.Printf("Session error: %v", err)
			return c.Redirect(http.StatusFound, "/admin/login")
		}

		if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
			return c.Redirect(http.StatusFound, "/admin/login")
		}
		return next(c)
	}
}

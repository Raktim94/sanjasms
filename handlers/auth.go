package handlers

import (
	"net/http"

	"odoo-lite-cms/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckFirstRun(c echo.Context) error {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		return c.Redirect(http.StatusFound, "/admin/setup")
	}
	// If user goes to /admin/setup but users exist, redirect to login
	if c.Path() == "/admin/setup" && count > 0 {
		return c.Redirect(http.StatusFound, "/admin/login")
	}
	return nil
}

func SetupPage(c echo.Context) error {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)

	// If it's the first run, always allow signup
	if count == 0 {
		return c.Render(http.StatusOK, "signup.html", nil)
	}

	// If not first run, check setting
	settings := GetSettingsMap()
	if settings["allow_public_signup"] != "true" {
		return c.Redirect(http.StatusFound, "/admin/login")
	}

	return c.Render(http.StatusOK, "signup.html", nil)
}

func SetupService(c echo.Context) error {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		settings := GetSettingsMap()
		if settings["allow_public_signup"] != "true" {
			return c.String(http.StatusForbidden, "Public signups are disabled")
		}
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.String(http.StatusBadRequest, "Email and password required")
	}

	hash, _ := HashPassword(password)
	user := models.User{Email: email, PasswordHash: hash, IsAdmin: true}
	models.DB.Create(&user)

	return c.Redirect(http.StatusFound, "/admin/login")
}

func LoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func LoginService(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	result := models.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{"Error": "Invalid credentials"})
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{"Error": "Invalid credentials"})
	}

	// Set session
	sess, _ := session.Get("admin_session", c)
	sess.Values["authenticated"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

func Logout(c echo.Context) error {
	sess, _ := session.Get("admin_session", c)
	sess.Values["authenticated"] = false
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/admin/login")
}

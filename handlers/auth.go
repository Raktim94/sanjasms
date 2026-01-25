package handlers

import (
	"net/http"
	"time"

	"odoo-lite-cms/models"

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
	if err := CheckFirstRun(c); err != nil {
		// If CheckFirstRun redirected, we shouldn't be here unless it returned nil (meaning no redirect)
		// But CheckFirstRun logic is: if count == 0, redirect to setup.
		// Wait, if we are AT /admin/setup, CheckFirstRun loop?
		// Let's refine the logic in routes, not here.
	}
	// Just check if users exist
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return c.Redirect(http.StatusFound, "/admin/login")
	}

	return c.Render(http.StatusOK, "signup.html", nil)
}

func SetupService(c echo.Context) error {
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

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "admin_session"
	cookie.Value = "logged_in" // In real app, use a secure token/session ID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/" // Important!
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "admin_session"
	cookie.MaxAge = -1
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/admin/login")
}

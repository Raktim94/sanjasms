package handlers

import (
	"net/http"
	"odoo-lite-cms/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func UserList(c echo.Context) error {
	var users []models.User
	models.DB.Find(&users)
	return c.Render(http.StatusOK, "users_list.html", map[string]interface{}{
		"Title":  "User Management",
		"Users":  users,
		"Active": "users",
	})
}

func UserNew(c echo.Context) error {
	return c.Render(http.StatusOK, "users_form.html", map[string]interface{}{
		"Title":  "Add New User",
		"Active": "users",
		"csrf":   c.Get("csrf"),
	})
}

func UserCreate(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	isAdmin := c.FormValue("is_admin") == "on"

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Email:        email,
		PasswordHash: string(hash),
		IsAdmin:      isAdmin,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/admin/users")
}

func UserDelete(c echo.Context) error {
	id := c.Param("id")
	// Prevent deleting last admin if needed
	models.DB.Delete(&models.User{}, id)
	return c.Redirect(http.StatusSeeOther, "/admin/users")
}

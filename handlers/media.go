package handlers

import (
	"fmt"
	"io"
	"net/http"
	"odoo-lite-cms/models"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func MediaList(c echo.Context) error {
	var files []models.File
	models.DB.Order("created_at desc").Find(&files)
	return c.Render(http.StatusOK, "media_list.html", map[string]interface{}{
		"Title":  "Media Manager",
		"Active": "media",
		"Files":  files,
	})
}

func MediaUpload(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	// Ensure filename is unique or just overwrite? Safe defaults -> Unique
	// Using timestamp prefix
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	dstPath := filepath.Join("static", "uploads", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Save to DB
	dbFile := models.File{
		Filename: file.Filename, // Original name
		Path:     "/uploads/" + filename,
		MimeType: file.Header.Get("Content-Type"),
		Size:     file.Size,
	}
	models.DB.Create(&dbFile)

	return c.Redirect(http.StatusFound, "/admin/media")
}

func MediaDelete(c echo.Context) error {
	id := c.Param("id")
	var file models.File
	if err := models.DB.First(&file, id).Error; err != nil {
		return c.String(http.StatusNotFound, "File not found")
	}

	// Delete from disk
	relPath := file.Path
	if len(relPath) > 0 && relPath[0] == '/' {
		relPath = relPath[1:]
	}
	fullPath := filepath.Join("static", relPath)

	os.Remove(fullPath) // Ignore error if not exists

	models.DB.Delete(&file)
	return c.Redirect(http.StatusFound, "/admin/media")
}

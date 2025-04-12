package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"backend/internal/handlers"
)

func Images(e *echo.Echo, db *gorm.DB) {
	e.GET("/images", handlers.GetImages(db))
	e.DELETE("/images/:id", handlers.DeleteImage(db))
} 
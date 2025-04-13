package handlers

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"backend/internal/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ImageResponse struct {
	ID        int       `json:"id"`
	ImageData string    `json:"image_data"`
	CreatedAt time.Time `json:"created_at"`
}

func GetImages(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		images, err := model.GetImages(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch images"})
		}

		response := make([]ImageResponse, len(images))
		for i, img := range images {
			response[i] = ImageResponse{
				ID:        img.ID,
				ImageData: base64.StdEncoding.EncodeToString(img.ImageData),
				CreatedAt: img.CreatedAt,
			}
		}

		return c.JSON(http.StatusOK, response)
	}
}

func DeleteImage(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid image ID"})
		}

		if err := model.DeleteImage(db, id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete image"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Image deleted successfully"})
	}
}

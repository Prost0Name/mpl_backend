package routes

import (
	"backend/internal/config"
	"backend/internal/handlers"

	"github.com/labstack/echo/v4"
)

func Users(e *echo.Echo, cfg *config.Config) {
	e.POST("/reg", func(c echo.Context) error {
		return handlers.Register(c, cfg)
	})
	e.POST("/auth", handlers.Login)
	e.GET("/confirm", handlers.ConfirmRegistration)
	// e.GET()
}

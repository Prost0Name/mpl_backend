package app

import (
	"backend/internal/app/middleware"
	"backend/internal/app/routes"
	"backend/internal/config"
	"backend/internal/model"
	"log"

	"github.com/labstack/echo/v4"
)

func New(cfg *config.Config) {
	e := echo.New()
	middleware.CORS(e)

	if err := model.InitDatabase(cfg.DSN); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	routes.Users(e, cfg)
	routes.Images(e, model.DB)

	if err := e.Start(":" + cfg.APP.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

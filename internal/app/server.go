package app

import (
	"backend/internal/app/middleware"
	"backend/internal/app/routes"
	"backend/internal/config"
	"backend/internal/model"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}
)

func handleWebSocket(c echo.Context, db *gorm.DB) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return err
	}
	defer ws.Close()

	log.Printf("New WebSocket connection established")

	for {
		endTime := time.Now()
		startTime := endTime.Add(-1 * time.Hour)

		log.Printf("Fetching data from %v to %v", startTime, endTime)

		data, err := model.CountProductionByInterval(db, startTime, endTime, "minute")
		if err != nil {
			log.Printf("Error getting production data: %v", err)
			// Отправляем пустой массив вместо null
			emptyData := []struct {
				Interval time.Time `json:"interval"`
				Count    int       `json:"count"`
			}{}
			if err := ws.WriteJSON(emptyData); err != nil {
				log.Printf("Error writing empty data to WebSocket: %v", err)
				return err
			}
			continue
		}

		log.Printf("Sending data: %+v", data)

		if err := ws.WriteJSON(data); err != nil {
			log.Printf("Error writing to WebSocket: %v", err)
			return err
		}

		time.Sleep(5 * time.Second)
	}
}

func New(cfg *config.Config) {
	e := echo.New()

	// Используем существующий CORS middleware
	middleware.CORS(e)

	if err := model.InitDatabase(cfg.DSN); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// Добавляем тестовые данные
	if err := model.AddTestData(model.DB); err != nil {
		log.Printf("Error adding test data: %v", err)
	}

	routes.Users(e, cfg)
	routes.Images(e, model.DB)

	// WebSocket endpoint
	e.GET("/ws/production", func(c echo.Context) error {
		return handleWebSocket(c, model.DB)
	})

	log.Printf("Server starting on port %s", cfg.APP.Port)
	if err := e.Start(":" + cfg.APP.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

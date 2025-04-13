package websocket

import (
	"log"
	"time"

	"backend/internal/model"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ProductionUpdate struct {
	Interval string    `json:"interval"`
	Count    int       `json:"count"`
	Time     time.Time `json:"time"`
}

func HandleProductionSocket(db *gorm.DB, conn *websocket.Conn) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		endTime := time.Now()
		startTime := endTime.Add(-1 * time.Hour)

		data, err := model.CountProductionByInterval(db, startTime, endTime, "5 minute")
		if err != nil {
			log.Printf("Error getting production data: %v", err)
			continue
		}

		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Error sending data through websocket: %v", err)
			return
		}
	}
}

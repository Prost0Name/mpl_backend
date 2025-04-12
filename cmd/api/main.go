package main

import (
	"backend/internal/app"
	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/redis"
	"backend/utils"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	redis.InitRedis()
	fmt.Println(cfg.JwtSecret)

	model.InitDatabase(cfg.DSN)
	
	// Start the email worker
	utils.StartEmailWorker()

	app.New(cfg)
}

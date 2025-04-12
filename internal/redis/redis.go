package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis() {
	ctx := context.Background()                                              
	opt, err := redis.ParseURL("redis://default:reg2025@87.242.100.33:6379") 
	if err != nil {
		log.Fatalf("Ошибка парсинга URL: %v", err)
	}

	Client = redis.NewClient(opt)

	// Тестируем соединение
	pong, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	log.Println("Успешное подключение:", pong)
}

func SetUserData(token string, data interface{}, ttl int) error {
	ctx := context.Background()
	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка при сериализации данных: %v", err)
		return err
	}
	// Устанавливаем данные в Redis с временем жизни, указанным в конфигурации
	err = Client.Set(ctx, token, jsonData, time.Duration(ttl)*time.Second).Err() // Сохраняем JSON-строку
	if err != nil {
		log.Printf("Ошибка при сохранении данных в Redis: %v", err) // Логируем ошибку
	}
	return err
}

func GetUserData(token string) (string, error) {
	ctx := context.Background()
	return Client.Get(ctx, token).Result()
}

package model

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type ProductionTime struct {
	ID          int       `gorm:"primaryKey"`
	Timestamp   time.Time `gorm:"column:time"`
	IsDefective bool      `gorm:"column:is_defective;default:false"`
}

func (ProductionTime) TableName() string {
	return "production_time"
}

func GetProductionTimes(db *gorm.DB, startTime, endTime time.Time) ([]ProductionTime, error) {
	var times []ProductionTime
	if err := db.Where("time BETWEEN ? AND ?", startTime, endTime).Find(&times).Error; err != nil {
		return nil, err
	}
	return times, nil
}

func CountProductionByInterval(db *gorm.DB, startTime, endTime time.Time, interval string) ([]struct {
	Interval     time.Time `json:"interval"`
	TotalCount   int       `json:"total_count"`   // Общее количество
	QualityCount int       `json:"quality_count"` // Количество качественных
}, error) {
	var result []struct {
		Interval     time.Time `json:"interval"`
		TotalCount   int       `json:"total_count"`
		QualityCount int       `json:"quality_count"`
	}

	query := `
		SELECT 
			date_trunc('minute', time) as interval,
			count(*) as total_count,
			count(CASE WHEN is_defective = false THEN 1 END) as quality_count
		FROM production_time
		WHERE time BETWEEN $1 AND $2
		GROUP BY interval
		ORDER BY interval
	`

	if err := db.Raw(query, startTime, endTime).Scan(&result).Error; err != nil {
		log.Printf("SQL Error: %v", err)
		return nil, err
	}

	log.Printf("Query result: %+v", result)

	if result == nil {
		result = make([]struct {
			Interval     time.Time `json:"interval"`
			TotalCount   int       `json:"total_count"`
			QualityCount int       `json:"quality_count"`
		}, 0)
	}

	return result, nil
}

func CreateProductionTime(db *gorm.DB, timestamp time.Time, isDefective bool) (*ProductionTime, error) {
	production := &ProductionTime{
		Timestamp:   timestamp,
		IsDefective: isDefective,
	}
	if err := db.Create(production).Error; err != nil {
		return nil, err
	}
	return production, nil
}

// AddTestData добавляет тестовые данные в базу
func AddTestData(db *gorm.DB) error {
	// Очищаем таблицу
	if err := db.Exec("DELETE FROM production_time").Error; err != nil {
		return err
	}

	// Добавляем тестовые данные за последний час
	now := time.Now()
	for i := 0; i < 12; i++ { // 12 записей с интервалом в 5 минут
		timestamp := now.Add(time.Duration(-i*5) * time.Minute)
		// Добавляем несколько записей для каждого интервала
		for j := 0; j < i+1; j++ {
			// Каждая третья деталь будет дефектной для тестовых данных
			isDefective := j%3 == 0
			_, err := CreateProductionTime(db, timestamp, isDefective)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

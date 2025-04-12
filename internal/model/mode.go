package model

import (
	"backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBUser struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
}

func (DBUser) TableName() string {
	return "users"
}

var DB *gorm.DB

func InitDatabase(dsn config.DSNConfig) error {
	var err error
	dsnString := "host=" + dsn.Host + " user=" + dsn.User + " password=" + dsn.Pass + " dbname=" + dsn.DBName + " port=" + dsn.Port + " sslmode=" + dsn.SSLMode
	DB, err = gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		return err
	}

	DB.AutoMigrate(&DBUser{})
	return nil
}

package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string     `yaml:"env" env-default:"local"`
	APP       APPConfig  `yaml:"app"`
	TokenTTL  int        `yaml:"token_ttl" env-default:"300"`
	JwtSecret string     `yaml:"JwtSecret"`
	DSN       DSNConfig  `yaml:"dsn"`
	SMTP      SMTPConfig `yaml:"smtp"`
}

type APPConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DSNConfig struct {
	Host    string `yaml:"host"`
	User    string `yaml:"user"`
	Pass    string `yaml:"password"`
	DBName  string `yaml:"dbname"`
	Port    string `yaml:"port"`
	SSLMode string `yaml:"sslmode"`
}

type SMTPConfig struct {
	From     string `yaml:"from"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func MustLoad() *Config {
	configPath := "config/local.yaml"
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

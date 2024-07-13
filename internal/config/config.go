package config

import (
	"log"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database struct {
		Host        string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
		Port        string `yaml:"port" env:"DB_PORT" env-description:"Database port"`
		Username    string `yaml:"username" env:"DB_USER" env-description:"Database user name"`
		Password    string `env:"DB_PASSWORD" env-description:"Database user password"`
		Name        string `yaml:"db-name" env:"DB_NAME" env-description:"Database name"`
		Connections int    `yaml:"connections" env:"DB_CONNECTIONS" env-description:"Total number of database connections"`
	} `yaml:"database"`
	Server struct {
		Host        string        `yaml:"host" env:"SRV_HOST,HOST" env-description:"Server host" env-default:"localhost"`
		Port        string        `yaml:"port" env:"SRV_PORT,PORT" env-description:"Server port" env-default:"8080"`
		Address     string        `yaml:"address" env-default:"localhost:8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	} `yaml:"server"`
	Greeting string `env:"GREETING" env-description:"Greeting phrase" env-default:"Hello!"`
}


func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
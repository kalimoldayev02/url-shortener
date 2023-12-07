package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yml:"env" env:"ENV" envDefault:"development"`
	HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Host        string        `yml:"host" env:"HOST"`
	Port        string        `yml:"port" env:"PORT"`
	Timeout     time.Duration `yml:"timeout" env:"TIMEOUT"`
	IdleTimeout time.Duration `yml:"idle_timeout" env:"IDLE_TIMEOUT"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error load .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cnf Config
	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("connot read config: %s", err)
	}

	return &cnf
}

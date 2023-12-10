package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

type Config struct {
	Env        string `yml:"env" env:"ENV" envDefault:"development"`
	HttpServer `yaml:"http_server"`
	DataBase   `yaml:"database"`
}

type HttpServer struct {
	Host        string        `yml:"host" env:"HOST"`
	Port        string        `yml:"port" env:"PORT"`
	Timeout     time.Duration `yml:"timeout" env:"TIMEOUT"`
	IdleTimeout time.Duration `yml:"idle_timeout" env:"IDLE_TIMEOUT"`
}

type DataBase struct {
	Host     string `yml:"host" env:"HOST"`
	Port     string `yml:"port" env:"PORT"`
	DbUser   string `yml:"user" env:"USER"`
	Password string `yml:"password" env:"PASSWORD"`
	Name     string `yml:"name" env:"NAME"`
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

func (cfg *Config) GetStoragePath() string {
	fmt.Println(cfg.DataBase.DbUser, cfg.DataBase) // TODO: надо исправить (чтобы через конфиг работал)
	return fmt.Sprintf("host=%s port=%s user=postgres password=%s dbname=%s sslmode=disable",
		cfg.DataBase.Host, cfg.DataBase.Port, cfg.DataBase.Password, cfg.DataBase.Name)
}

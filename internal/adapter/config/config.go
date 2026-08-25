package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App      *App
		HTTP     *HTTP
		DB       *DB
		Cache    *Cache
		Rabbitmq *Rabbitmq
	}

	App struct {
		Name string
		Env  string
	}

	HTTP struct {
		Host           string
		Port           string
		AllowedOrigins string
	}

	DB struct {
		Name     string
		Host     string
		Port     string
		User     string
		Password string
	}

	Cache struct {
		Host     string
		Port     string
		Password string
	}

	Rabbitmq struct {
		Host     string
		Port     string
		User     string
		Password string
	}
)

func New() (*Container, error) {
	_ = godotenv.Load()

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	http := &HTTP{
		Host:           os.Getenv("HTTP_HOST"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	db := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	cache := &Cache{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	rabbitmq := &Rabbitmq{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}

	return &Container{
		App:      app,
		HTTP:     http,
		DB:       db,
		Cache:    cache,
		Rabbitmq: rabbitmq,
	}, nil
}

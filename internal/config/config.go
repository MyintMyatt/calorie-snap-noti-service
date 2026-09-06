package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	RabbitMQUrl    string
}

func LoadConfig() (*Config, error) {

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	fmt.Printf("ENV PROFILE: ==== %s ====\n", env)

	if err := godotenv.Load(".env." + env); err != nil {
		return nil, fmt.Errorf("error loading .env.%s file: %v", env, err)
	}

	cfg := &Config{
		AppEnv: os.Getenv("APP_ENV"),
		RabbitMQUrl:    os.Getenv("RABBITMQ_URL"),
	}

	if cfg.RabbitMQUrl == "" {
		return nil, fmt.Errorf("RABBITMQ_URL is required")
	}

	return cfg, nil
}
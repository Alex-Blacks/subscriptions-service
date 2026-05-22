package config

import "os"

type Config struct {
	AppPort string

	DatabaseURL string
}

func MustLoad() Config {
	return Config{
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

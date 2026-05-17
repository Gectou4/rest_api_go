package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	User string
	Pwd  string
	DSN  string
}

func LoadDBConfig() DBConfig {
	cfg := DBConfig{
		User: getEnv("DB_USER", "root"),
		Pwd:  getEnv("DB_PWD", ""),
		DSN:  getEnv("DB_DSN", "mysql:host=localhost;dbname=rest_api;charset=utf8"),
	}
	return cfg
}

func (c DBConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@%s", c.User, c.Pwd, c.DSN)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

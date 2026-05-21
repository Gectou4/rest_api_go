package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	User string
	Pwd  string
	Host string
	Port string
	Name string
	TLS  string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		User: getEnv("DB_USER", "root"),
		Pwd:  getEnv("DB_PWD", ""),
		Host: getEnv("DB_HOST", "127.0.0.1"),
		Port: getEnv("DB_PORT", "3306"),
		Name: getEnv("DB_NAME", "rest_api"),
		TLS:  getEnv("DB_TLS", ""),
	}
}

func (c DBConfig) ConnectionString() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.User, c.Pwd, c.Host, c.Port, c.Name,
	)
	if c.TLS != "" {
		dsn += "&tls=" + c.TLS
	}
	return dsn
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

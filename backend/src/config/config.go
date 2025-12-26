package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	MusicDir   string
	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "wsb"),
		DBPassword: getEnv("DB_PASSWORD", "admin1234"),
		DBName:     getEnv("DB_NAME", "myapp"),
		MusicDir:   getEnv("MUSIC_DIR", "./songs"),
		ServerPort: getEnv("SERVER_PORT", "8888"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

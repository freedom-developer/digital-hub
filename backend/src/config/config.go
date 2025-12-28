package config

import (
	"myapp/database"
	logger "myapp/log"
	"myapp/servers/music"
	"os"
)

type Config struct {
	// 数据库配置
	DbConfig database.DbConfig

	// 日志配置
	LogConfig logger.LogConfig

	// HTTP服务器配置
	SrvConfig struct {
		Addr string
		Port string
	}

	// 音乐配置
	MusicConfig music.MusicConfig
}

func LoadConfig() *Config {
	return &Config{
		DbConfig: database.DbConfig{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "3306"),
			DBUser:     getEnv("DB_USER", "wsb"),
			DBPassword: getEnv("DB_PASSWORD", "admin1234"),
			DBName:     getEnv("DB_NAME", "myapp"),
			MusicDir:   getEnv("MUSIC_DIR", "./songs"),
			ServerPort: getEnv("SERVER_PORT", "8888"),
		},

		LogConfig: logger.LogConfig{
			LogDir: getEnv("LOG_DIR", "./logs"),
			LogLevel: func() int {
				level := getEnv("LOG_LEVEL", "1")
				switch level {
				case "1":
					return 1
				case "2":
					return 2
				case "3":
					return 3
				case "4":
					return 4
				case "5":
					return 5
				default:
					return 2
				}
			}(),
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},

		SrvConfig: struct {
			Addr string
			Port string
		}{
			"0.0.0.0",
			"8888",
		},

		MusicConfig: music.MusicConfig{
			MusicDir: getEnv("MUSIC_DIR", "./songs"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

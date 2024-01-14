package config

import (
	"os"
)

type Config struct {
	App    *AppConfig
	Server *ServerConfig
	DB     *DBConfig
}

func New() *Config {
	return &Config{
		App:    NewAppConfig(),
		Server: NewServerConfig(),
		DB:     NewDBConfig(),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

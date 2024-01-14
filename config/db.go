package config

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     getEnv("DB_HOST", ""),
		Port:     getEnv("DB_PORT", ""),
		Name:     getEnv("DB_DATABASE", ""),
		User:     getEnv("DB_USERNAME", ""),
		Password: getEnv("DB_PASSWORD", ""),
	}
}

package config

type Environment string

const (
	Local       Environment = "local"
	Production  Environment = "production"
	Development Environment = "development"
	Test        Environment = "test"
)

type AppConfig struct {
	Environment string
	Name        string
	Version     string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Environment: getEnv("APP_ENV", string(Production)),
		Name:        getEnv("APP_NAME", ""),
		Version:     getEnv("APP_VERSION", ""),
	}
}

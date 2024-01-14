package config

type ServerConfig struct {
	HttpAddr string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		HttpAddr: getEnv("HTTP_SERVER_ADDR", ":8080"),
	}
}

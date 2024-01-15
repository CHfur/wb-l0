package config

type NatsConfig struct {
	ClusterId string
	ClientId  string
	Url       string
	Channel   string
}

func NewNatsConfig() *NatsConfig {
	return &NatsConfig{
		ClusterId: getEnv("NATS_CLUSTER_ID", ""),
		ClientId:  getEnv("NATS_CLIENT_ID", ""),
		Url:       getEnv("NATS_URL", ""),
		Channel:   getEnv("NATS_CHANNEL", ""),
	}
}

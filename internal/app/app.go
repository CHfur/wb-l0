package app

import (
	"go.uber.org/zap"
	"order-service/config"
	httpapp "order-service/internal/app/http"
	natsapp "order-service/internal/app/nats"
	"order-service/internal/services"
	"order-service/internal/storage/pgsql"
)

type App struct {
	HTTPServer   *httpapp.App
	NATSListener *natsapp.App
}

func NewApp(logger *zap.Logger, httpAddr string, databaseConfig *config.DBConfig, natsConfig *config.NatsConfig) *App {
	storage, err := pgsql.New(databaseConfig)
	if err != nil {
		panic(err)
	}

	cache := services.LoadOrdersFromDbToCache(logger, storage)
	service := services.NewService(logger, storage, storage, cache)

	httpApp := httpapp.NewApp(logger, service, httpAddr)
	natsApp := natsapp.NewApp(logger, service, natsConfig.ClusterId, natsConfig.ClientId, natsConfig.Url, natsConfig.Channel)

	return &App{HTTPServer: httpApp, NATSListener: natsApp}
}

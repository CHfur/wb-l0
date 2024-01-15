package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"order-service/config"
	"order-service/internal/app"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.New()

	logger, err := createLogger(conf.App.Name, conf.App.Version, "debug")
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	defer cancel()

	application := app.NewApp(logger, conf.Server.HttpAddr, conf.DB, conf.Nats)

	go application.HTTPServer.MustRun()

	<-ctx.Done()

	logger.Info("got signal to shutdown")

	application.HTTPServer.Stop()

	err = application.NATSListener.Stop()
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("application stopped")
}

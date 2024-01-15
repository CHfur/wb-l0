package httpapp

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"order-service/internal/services"
	"time"
)

const DefaultReadTimeout = 5 * time.Second

type App struct {
	logger     *zap.Logger
	httpServer *http.Server
}

func NewApp(logger *zap.Logger, service *services.Service, httpAddr string) *App {
	controller := NewController(service, logger)

	server := &http.Server{
		Addr:        httpAddr,
		ReadTimeout: DefaultReadTimeout,
		Handler:     CreateRouter(controller),
	}

	return &App{logger: logger, httpServer: server}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.logger.With(
		zap.String("op", op),
	)

	log.Info("HTTP server is running", zap.String("addr", a.httpServer.Addr))

	return a.httpServer.ListenAndServe()
}

func (a *App) Stop() {
	const op = "httpapp.Stop"

	a.logger.With(zap.String("op", op)).Info("stopping HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("cannot shutdown HTTP server", zap.Error(err))
		if err := a.httpServer.Close(); err != nil {
			a.logger.Error("cannot close HTTP server", zap.Error(err))
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}

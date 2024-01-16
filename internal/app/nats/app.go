package natsapp

import (
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"go.uber.org/zap"
	"order-service/internal/services"
)

type App struct {
	logger *zap.Logger
	conn   stan.Conn
	sub    stan.Subscription
	cancel context.CancelFunc
}

func NewApp(logger *zap.Logger, service *services.Service, clusterId, clientId, url, channel string) *App {
	const op = "natsapp.NewApp"

	log := logger.With(zap.String("op", op))

	handler := NewHandler(service, logger)

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(url),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Error(fmt.Sprintf("Connection lost, reason: %v", reason))
		}))
	if err != nil {
		log.Error(fmt.Sprintf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url))
	}
	log.Info(fmt.Sprintf("Connected to nats-streaming %s clusterID: [%s] clientID: [%s]", url, clusterId, clientId))

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	sub, err := sc.QueueSubscribe(channel, "", handler.SaveOrder, startOpt, stan.DurableName(""))
	if err != nil {
		sc.Close()
		panic(err)
	}

	log.Info(fmt.Sprintf("Nats-streaming listening on [%s], clientID=[%s]", channel, clientId))

	ctx, cancel := context.WithCancel(context.Background())
	handler.StartHandle(ctx)
	log.Info("Nats handler started")

	return &App{logger: logger, conn: sc, sub: sub, cancel: cancel}
}

func (a *App) Stop() error {
	const op = "natsapp.Stop"

	a.logger.With(zap.String("op", op)).Info("stopping NATS stream handle")

	a.cancel()
	errSub := a.sub.Close()
	errConn := a.conn.Close()

	if errConn != nil {
		return fmt.Errorf("%s: %w", op, errConn)
	}

	if errSub != nil {
		return fmt.Errorf("%s: %w", op, errSub)
	}

	return nil
}

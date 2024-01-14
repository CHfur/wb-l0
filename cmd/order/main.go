package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"go.uber.org/zap"
	"log"
	net_http "net/http"
	"order-service/config"
	"order-service/http"
	"order-service/internal/domain/models"
	"order-service/internal/services"
	"order-service/internal/storage/pgsql"
	"os/signal"
	"syscall"
	"time"
)

func printMsg(m *stan.Msg, i int) {
	//fmt.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, m)
}

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

	storage, err := pgsql.New(conf.DB)
	if err != nil {
		panic(err)
	}

	service := services.NewService(logger, storage, storage)
	service.LoadOrdersFromDbToCache()

	controller := http.NewController(service, logger)

	conn, sub := startStanListening(service, logger)

	errors := make(chan error, 1)

	server := net_http.Server{
		Addr:        conf.Server.HttpAddr,
		ReadTimeout: 5 * time.Second,
		Handler:     http.CreateRouter(controller),
	}
	go func() {
		logger.Info("server started", zap.String("port", conf.Server.HttpAddr))
		errors <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("got signal to shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn.Close()
		sub.Close()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("cannot shutdown server", zap.Error(err))
			if err := server.Close(); err != nil {
				logger.Error("cannot close server", zap.Error(err))
			}
		}

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

	case err = <-errors:
		if err != nil {
			logger.Fatal("cannot start service", zap.Error(err))
		}
	}
}

func startStanListening(service *services.Service, logger *zap.Logger) (stan.Conn, stan.Subscription) {
	clusterID := "test-cluster"
	clientID := "test-client-1"
	showTime := false
	var startSeq uint64 = 0
	startDelta := ""
	deliverAll := false
	deliverLast := false
	durable := ""
	qgroup := ""
	URL := "localhost:4222"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	logger.Info(fmt.Sprintf("Connected to %s clusterID: [%s] clientID: [%s]", URL, clusterID, clientID))

	subj, i := "orders", 0

	mcb := func(msg *stan.Msg) {
		i++
		printMsg(msg, i)

		var order models.Order

		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			logger.Info("Wrong data format")
		}

		err = service.SaveOrder(order)
		if err != nil {
			logger.Info(err.Error())
		}
	}

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	if startSeq != 0 {
		startOpt = stan.StartAtSequence(startSeq)
	} else if deliverLast {
		startOpt = stan.StartWithLastReceived()
	} else if deliverAll {
		log.Print("subscribing with DeliverAllAvailable")
		startOpt = stan.DeliverAllAvailable()
	} else if startDelta != "" {
		ago, err := time.ParseDuration(startDelta)
		if err != nil {
			sc.Close()
			log.Fatal(err)
		}
		startOpt = stan.StartAtTimeDelta(ago)
	}

	sub, err := sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Listening on [%s], clientID=[%s]", subj, clientID))

	if showTime {
		log.SetFlags(log.LstdFlags)
	}

	return sc, sub
}

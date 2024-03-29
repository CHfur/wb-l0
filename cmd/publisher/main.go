package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"order-service/config"
	"sync"
)

const ProduceOrdersCount = 5000

func main() {
	conf := config.New()

	logger, err := createLogger(conf.App.Name, conf.App.Version, "debug")
	if err != nil {
		panic(err)
	}

	clusterID := conf.Nats.ClusterId
	clientID := "test-client-2"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(conf.Nats.Url))
	if err != nil {
		panic(err)
	}
	defer sc.Close()

	var wg sync.WaitGroup

	ackHandler := func(ackedNuid string, err error) {
		defer wg.Done()
		if err != nil {
			logger.Info(fmt.Sprintf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error()))
		} else {
			logger.Info(fmt.Sprintf("Received ack for msg id %s\n", ackedNuid))
		}
	}

	for i := 0; i < ProduceOrdersCount; i++ {
		uid := RandStringRunes(20)
		wg.Add(1)

		_, err = sc.PublishAsync(conf.Nats.Channel, getMessage(uid), ackHandler)
		if err != nil {
			panic(err)
		}
	}

	wg.Wait()

	logger.Info("Order successfully produced")
}

func createLogger(name, version, level string) (*zap.Logger, error) {
	zl := zap.InfoLevel
	if err := zl.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	opts := zap.NewProductionConfig()
	opts.Level = zap.NewAtomicLevelAt(zl)
	opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	opts.Encoding = "console"
	opts.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := opts.Build()
	if err != nil {
		return nil, err
	}

	return logger.With(
		zap.String("name", name),
		zap.String("version", version),
	), nil
}

var letterRunes = []rune("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getMessage(uid string) []byte {
	return []byte("{\n  \"order_uid\": \"" + uid + "\",\n  \"track_number\": \"WBILMTESTTRACK\",\n  \"entry\": \"WBIL\",\n  \"delivery\": {\n    \"name\": \"Test Testov\",\n    \"phone\": \"+9720000000\",\n    \"zip\": \"2639809\",\n    \"city\": \"Kiryat Mozkin\",\n    \"address\": \"Ploshad Mira 15\",\n    \"region\": \"Kraiot\",\n    \"email\": \"test@gmail.com\"\n  },\n  \"payment\": {\n    \"transaction\": \"b563feb7b2b84b6test\",\n    \"request_id\": \"\",\n    \"currency\": \"USD\",\n    \"provider\": \"wbpay\",\n    \"amount\": 1817,\n    \"payment_dt\": 1637907727,\n    \"bank\": \"alpha\",\n    \"delivery_cost\": 1500,\n    \"goods_total\": 317,\n    \"custom_fee\": 0\n  },\n  \"items\": [\n    {\n      \"chrt_id\": 9934930,\n      \"track_number\": \"WBILMTESTTRACK\",\n      \"price\": 453,\n      \"rid\": \"ab4219087a764ae0btest\",\n      \"name\": \"Mascaras\",\n      \"sale\": 30,\n      \"size\": \"0\",\n      \"total_price\": 317,\n      \"nm_id\": 2389212,\n      \"brand\": \"Vivienne Sabo\",\n      \"status\": 202\n    }\n  ],\n  \"locale\": \"en\",\n  \"internal_signature\": \"\",\n  \"customer_id\": \"test\",\n  \"delivery_service\": \"meest\",\n  \"shardkey\": \"9\",\n  \"sm_id\": 99,\n  \"date_created\": \"2021-11-26T06:22:19Z\",\n  \"oof_shard\": \"1\"\n}")
}

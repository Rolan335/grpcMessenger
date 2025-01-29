package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Rolan335/grpcMessenger/kafkaconsumer/internal/kafka"
	"github.com/Rolan335/grpcMessenger/kafkaconsumer/internal/repository/inmemory"
	"github.com/Rolan335/grpcMessenger/kafkaconsumer/internal/webhook"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("cannot load .env file: " + err.Error())
	}
	kafka.Init(os.Getenv("KAFKA_BROKER_0"))

	webhookMaxRetries, err := strconv.Atoi(os.Getenv("WEBHOOK_MAX_RETRIES"))
	if err != nil {
		panic("failed to parse WEBHOOK_MAX_RETRIES .env:" + err.Error())
	}
	webhookTimeout, err := strconv.Atoi(os.Getenv("WEBHOOK_TIMEOUT"))
	if err != nil {
		panic("failed to parse WEBHOOK_TIMEOUT .env:" + err.Error())
	}
	wh := webhook.NewCaller(
		os.Getenv("WEBHOOK_METHOD"),
		os.Getenv("WEBHOOK_URL"),
		webhookMaxRetries,
		webhookTimeout,
	)

	//init storage to store processed kafka messages
	storage := inmemory.NewStorageTTL(time.Second * 20)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go kafka.ConsumeCreateChat(wh, storage)
	<-ctx.Done()
	kafka.Close()
}

//go:generate protoc --grpc-gateway_out=../pkg/proto --proto_path=../pkg/proto --go_out=../pkg/proto --go-grpc_out=../pkg/proto --go_opt=paths=source_relative ../pkg/proto/messenger.proto
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Rolan335/grpcMessenger/server/internal/app"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/metric"

	"github.com/joho/godotenv"
)

func main() {
	//load .env
	err := godotenv.Load(".env.local")
	if err != nil {
		panic("cannot load .env file: " + err.Error())
	}

	maxChatSize, err := strconv.Atoi(os.Getenv("APP_MAXCHATSIZE"))
	if err != nil {
		panic("failed to parse maxChatSize .env:" + err.Error())
	}
	maxChats, err := strconv.Atoi(os.Getenv("APP_MAXCHATS"))
	if err != nil {
		panic("failed to parse maxChats .env:" + err.Error())
	}
	serverConfig := config.MustConfigInit(
		os.Getenv("APP_ADDRESS"),
		os.Getenv("APP_PORTGRPC"),
		os.Getenv("APP_PORTHTTP"),
		os.Getenv("APP_ENV"),
		maxChatSize,
		maxChats,
		os.Getenv("APP_DB"),
	)

	//initializing logger
	logger := logger.Init(serverConfig.Env, os.Stdout)

	//registering metrics
	metric.MustInit()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	app := app.NewServiceServer(serverConfig, logger)
	go app.MustStartGRPC()
	go app.MustStartHTTP(ctx)
	<-ctx.Done()
	app.GracefulStop()
}

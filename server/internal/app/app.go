package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/app/health"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/controller"
	"github.com/Rolan335/grpcMessenger/server/internal/controller/interceptors"
	"github.com/Rolan335/grpcMessenger/server/internal/kafka"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/inmemory"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/postgres"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/redis"
	"github.com/Rolan335/grpcMessenger/server/internal/service/messenger"
	"github.com/Rolan335/grpcMessenger/server/pkg/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServiceServer struct {
	config        config.ServiceCfg
	grpcServer    *grpc.Server
	httpServerMux *runtime.ServeMux
	httpServer    *http.Server
	storage       messenger.Storage
}

func NewServiceServer(config config.ServiceCfg) *ServiceServer {
	//kafka initialization
	kafka.Init(os.Getenv("KAFKA_BROKER_0"))

	//making interceptors chain
	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(
		interceptors.Metric,
		interceptors.Log,
	)

	serverOptions := []grpc.ServerOption{chainUnaryInterceptor}
	//Make new grpc server instance
	s := grpc.NewServer(serverOptions...)

	db := storageInit(config.StorageType, config.MaxChats, config.MaxChatSize)

	return &ServiceServer{
		config:        config,
		grpcServer:    s,
		httpServerMux: runtime.NewServeMux(),
		httpServer:    &http.Server{Addr: config.PortHTTP, Handler: nil},
		storage:       db,
	}
}

func storageInit(storageType string, maxChats int, maxChatSize int) messenger.Storage {
	//initializing storage
	var db messenger.Storage
	switch storageType {
	case "postgres":
		port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		if err != nil {
			panic("failed to parse .env POSTGRES_PORT: " + err.Error())
		}
		freshstart, err := strconv.ParseBool(os.Getenv("POSTGRES_FLUSH"))
		if err != nil {
			panic("failed to parse .env POSTGRES_FLUSH: " + err.Error())
		}
		postgresConfig := postgres.Config{
			Host:       os.Getenv("POSTGRES_HOST"),
			User:       os.Getenv("POSTGRES_USER"),
			Password:   os.Getenv("POSTGRES_PASSWORD"),
			Dbname:     os.Getenv("POSTGRES_DBNAME"),
			Port:       port,
			FreshStart: freshstart,
		}
		timeout := 5
		db = postgres.NewStorage(postgresConfig, maxChats, maxChatSize, timeout, os.Getenv("MIGRATIONS_PATH"))
	case "redis":
		freshstart, err := strconv.ParseBool(os.Getenv("REDIS_FLUSH"))
		if err != nil {
			panic("failed to parse .env REDIS_FLUSH: " + err.Error())
		}
		redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			panic("failed to parse .env REDIS_DB: " + err.Error())
		}
		redisConfig := redis.Config{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDb,
			FlushAll: freshstart,
		}
		db = redis.NewStorage(redisConfig, maxChatSize, maxChats)
	default:
		db = inmemory.NewStorage(maxChatSize, maxChats)
	}
	return db
}

func (s *ServiceServer) MustStartGRPC() {
	//start tcp connection
	lis, err := net.Listen("tcp", s.config.PortGRPC)
	if err != nil {
		panic(err)
	}
	//Registering server
	proto.RegisterMessengerServiceServer(s.grpcServer, controller.NewServer(s.storage))
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}

func (s *ServiceServer) MustStartHTTP(ctx context.Context) {
	err := proto.RegisterMessengerServiceHandlerFromEndpoint(
		ctx,
		s.httpServerMux,
		s.config.Address+s.config.PortGRPC,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		panic(err)
	}
	//all requests go through mux handler except metrics, health, ready
	http.Handle("/", s.httpServerMux)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", health.HealthHandler)
	http.HandleFunc("/ready", health.ReadyHandler)
	go func() {
		if err := http.ListenAndServe(s.config.PortHTTP, nil); err != nil {
			panic("cannot start http server: " + err.Error())
		}
	}()
}

func (s *ServiceServer) GracefulStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Logger.ErrorContext(ctx, "failed to gracefully shutdown http: "+err.Error())
	}
	s.grpcServer.GracefulStop()
	redis.GracefulStop()
	postgres.GracefulStop()
	kafka.Close()
	logger.Logger.Info("gracefully shutdown")
}

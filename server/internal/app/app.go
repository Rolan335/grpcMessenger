package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/controller"
	"github.com/Rolan335/grpcMessenger/server/internal/controller/interceptors"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/postgres"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/redis"
	"github.com/Rolan335/grpcMessenger/server/pkg/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceServer struct {
	config        config.ServiceCfg
	grpcServer    *grpc.Server
	httpServerMux *runtime.ServeMux
	logger        logger.Logger
}

func NewServiceServer(config config.ServiceCfg, logger logger.Logger) *ServiceServer {
	//making interceptors chain
	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(interceptors.Metric, interceptors.Log(logger))
	//making keepalive rules
	timeout := grpc.ConnectionTimeout(5 * time.Second)
	serverOptions := []grpc.ServerOption{chainUnaryInterceptor, timeout}
	//Make new grpc server instance
	s := grpc.NewServer(serverOptions...)

	return &ServiceServer{
		config:        config,
		grpcServer:    s,
		httpServerMux: runtime.NewServeMux(),
		logger:        logger,
	}
}

func (s *ServiceServer) MustStartGrpc() {
	//start tcp connection
	lis, err := net.Listen("tcp", s.config.PortGrpc)
	if err != nil {
		panic(err)
	}
	//Registering server
	proto.RegisterMessengerServiceServer(s.grpcServer, controller.NewServer(s.config, s.logger))
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}

func (s *ServiceServer) MustStartHttp(ctx context.Context) {
	err := proto.RegisterMessengerServiceHandlerFromEndpoint(
		ctx,
		s.httpServerMux,
		s.config.Address+s.config.PortGrpc,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		panic(err)
	}
	//all requests go through mux handler except metrics
	http.Handle("/", s.httpServerMux)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(s.config.PortHttp, nil); err != nil {
		panic("cannot start http server: " + err.Error())
	}
}

func (a *ServiceServer) GracefulStop() {
	a.grpcServer.GracefulStop()
	redis.GracefulStop()
	postgres.GracefulStop()
}

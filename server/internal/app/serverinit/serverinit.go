package serverinit

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/controller"
	"github.com/Rolan335/grpcMessenger/server/internal/controller/interceptors"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/proto"
	"github.com/prometheus/client_golang/prometheus"

	"google.golang.org/grpc"
)

func Start(serverConfig config.ServiceCfg) {
	//Starting tcp connection
	lis, err := net.Listen("tcp", serverConfig.PortGrpc)
	if err != nil {
		fmt.Println(err)
	}
	
	//ctrl+c signal chan for graceful stop
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)

	//Registering metrics
	prometheus.MustRegister(interceptors.RequestsCounter)
	logger.Logger.Info("start prometheus metric at :8080/metrics")

	//making interceptors chain
	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(interceptors.Log, interceptors.Metric)

	//Starting grpc server and our service
	s := grpc.NewServer(chainUnaryInterceptor)
	proto.RegisterMessengerServiceServer(s, controller.NewServer(serverConfig))

	//started server concurrently to make graceful stop
	go func() {
		err = s.Serve(lis)
		if err != nil {
			logger.Logger.Error("server cannot start. Error: " + err.Error())
			panic(err)
		}
	}()
	logger.Logger.Info("started grpc server")

	//starting http router
	ctx, cancel := context.WithCancel(context.Background())
	go StartHttp(ctx, serverConfig)

	//graceful stop implementation
	<-stopSig
	cancel()
	s.GracefulStop()
	logger.Logger.Info("server stopped")
}

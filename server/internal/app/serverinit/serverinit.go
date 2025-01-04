package serverinit

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rolan335/grpcMessenger/proto"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/controller"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"google.golang.org/grpc"
)

func Start(serverConfig config.ServiceCfg) {
	//Logs are written to file, initializing logger
	f, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("cannot start logging " + err.Error())
	}
	defer f.Close()
	logger.LoggerInit(serverConfig.Env, f)

	//Starting tcp connection
	lis, err := net.Listen("tcp", serverConfig.Port)
	if err != nil {
		fmt.Println(err)
	}
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)

	//Starting grpc server and our service
	s := grpc.NewServer()
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

	//graceful stop implementation
	<-stopSig
	s.GracefulStop()
	logger.Logger.Info("server stopped")
}

package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/Rolan335/grpcMessenger/proto"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/controller"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"google.golang.org/grpc"
)

var serverConfig config.ServiceCfg

func main() {
	//Parsing flags to config
	port := flag.String("port", ":50051", "port of the app. Default [:50051]")
	maxChatSize := flag.Int("maxChatSize", 100, "maximum of messages that chat will store")
	maxChats := flag.Int("maxChats", 10, "maximum of chats that app will store")
	env := flag.String("env", "dev", "environment")
	flag.Parse()
	serverConfig = config.ServiceInit(*port, *env, *maxChatSize, *maxChats)

	//Logs are written to file, initializing logger
	f, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic("cannot start logging " + err.Error())
	}
	logger.LoggerInit(serverConfig.Env, f)

	//Starting tcp connection
	lis, err := net.Listen("tcp", serverConfig.Port)
	if err != nil {
		fmt.Println(err)
	}

	//Starting grpc server and our service
	s := grpc.NewServer()
	proto.RegisterMessengerServiceServer(s, controller.NewServer(serverConfig))
	logger.Logger.Info("started grpc server")
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}

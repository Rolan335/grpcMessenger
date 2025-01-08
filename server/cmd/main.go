//go:generate protoc --grpc-gateway_out=../../proto --proto_path=../../proto --go_out=../../proto --go-grpc_out=../../proto --go_opt=paths=source_relative ../../proto/messenger.proto
package main

import (
	"flag"
	"os"

	"github.com/Rolan335/grpcMessenger/server/internal/app/serverinit"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
)

func main() {
	//Parsing flags to config
	port := flag.String("port", ":50051", "port of the app")
	maxChatSize := flag.Int("maxchatsize", 100, "maximum of messages that chat will store")
	maxChats := flag.Int("maxchats", 10, "maximum of chats that app will store")
	env := flag.String("env", "dev", "environment")
	flag.Parse()
	serverConfig := config.ServiceInit("localhost", *port, ":8080", *env, *maxChatSize, *maxChats)

	//Logs are written to file, initializing logger
	f, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic("cannot start logging " + err.Error())
	}
	defer f.Close()
	logger.LoggerInit(serverConfig.Env, f)

	serverinit.Start(serverConfig)
}

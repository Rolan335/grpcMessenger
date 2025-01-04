package main

import (
	"flag"

	"github.com/Rolan335/grpcMessenger/server/internal/app/serverinit"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
)

var serverConfig config.ServiceCfg

func main() {
	//Parsing flags to config
	port := flag.String("port", ":50051", "port of the app")
	maxChatSize := flag.Int("maxchatsize", 100, "maximum of messages that chat will store")
	maxChats := flag.Int("maxchats", 10, "maximum of chats that app will store")
	env := flag.String("env", "dev", "environment")
	flag.Parse()
	serverConfig = config.ServiceInit(*port, *env, *maxChatSize, *maxChats)

	serverinit.Start(serverConfig)
}

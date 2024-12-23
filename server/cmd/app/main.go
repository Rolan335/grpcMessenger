package main

import (
	"flag"
	
	"github.com/Rolan335/grpcMessenger/server/internal/config"
)

func main() {
	port := flag.String("port", ":8080", "port of the app")
	maxChatSize := flag.Int("maxChatSize", 100, "maximum of messages that chat will store")
	maxChats := flag.Int("maxChats", 10, "maximum of chats that app will store")
	flag.Parse()
	_ = config.AppInit(*port, *maxChatSize, *maxChats)
}

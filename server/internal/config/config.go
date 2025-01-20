package config

import (
	"regexp"
	"strconv"
)

// Config that provided from flags at start
type ServiceCfg struct {
	Address     string
	PortGRPC    string
	PortHTTP    string
	Env         string
	MaxChatSize int
	MaxChats    int
	StorageType string
}

// Initializing Config and panic if couldn't
func MustConfigInit(address string, portGRPC string, portHTTP string, env string, maxChatSize int, maxChats int, storageType string) ServiceCfg {
	maxChatsAvailable := 1000
	maxChatSizeAvailable := 5000
	maxPort := 65535

	reg := regexp.MustCompile(`:\d{1,5}`)
	if portGRPC == portHTTP {
		panic("equal port for grpc and http")
	}
	portHTTPInt, _ := strconv.Atoi(portHTTP[1:])
	portGRPCInt, _ := strconv.Atoi(portGRPC[1:])

	if !reg.MatchString(portGRPC) || portGRPCInt > maxPort {
		panic("grpc port invalid")
	}

	if !reg.MatchString(portHTTP) || portHTTPInt > maxPort {
		panic("http port invalid")
	}

	if maxChats > maxChatsAvailable {
		maxChats = maxChatsAvailable
	}

	if maxChatSize > maxChatSizeAvailable {
		maxChatSize = maxChatSizeAvailable
	}

	return ServiceCfg{
		Address:     address,
		PortGRPC:    portGRPC,
		PortHTTP:    portHTTP,
		Env:         env,
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		StorageType: storageType,
	}
}

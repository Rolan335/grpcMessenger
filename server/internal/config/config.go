package config

import (
	"regexp"
	"strconv"
)

// Config that provided from flags at start
type ServiceCfg struct {
	Address     string
	PortGrpc    string
	PortHttp    string
	Env         string
	MaxChatSize int
	MaxChats    int
	StorageType string
}

// Initializing Config and panic if couldn't
func MustConfigInit(address string, portGrpc string, portHttp string, env string, maxChatSize int, maxChats int, storageType string) ServiceCfg {
	maxChatsAvailable := 1000
	maxChatSizeAvailable := 5000
	maxPort := 65535

	reg := regexp.MustCompile(`:\d{1,5}`)
	if portGrpc == portHttp {
		panic("equal port for grpc and http")
	}
	portHttpInt, _ := strconv.Atoi(portHttp[1:])
	portGrpcInt, _ := strconv.Atoi(portGrpc[1:])

	if !reg.MatchString(portGrpc) || portGrpcInt > maxPort {
		panic("grpc port invalid")
	}

	if !reg.MatchString(portHttp) || portHttpInt > maxPort {
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
		PortGrpc:    portGrpc,
		PortHttp:    portHttp,
		Env:         env,
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		StorageType: storageType,
	}
}

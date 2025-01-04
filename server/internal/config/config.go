package config

//Config that provided from flags at start
type ServiceCfg struct {
	Address     string
	PortGrpc    string
	PortHttp    string
	Env         string
	MaxChatSize int
	MaxChats    int
}

//Initializing Config
func ServiceInit(address string, portGrpc string, portHttp string, env string, maxChatSize int, maxChats int) ServiceCfg {
	return ServiceCfg{
		Address:     address,
		PortGrpc:    portGrpc,
		PortHttp:    portHttp,
		Env:         env,
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
	}
}

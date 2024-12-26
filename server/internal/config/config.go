package config

//Config that provided from flags at start
type ServiceCfg struct {
	Port        string
	Env         string
	MaxChatSize int
	MaxChats    int
}

//Initializing Config
func ServiceInit(port string, env string, maxChatSize int, maxChats int) ServiceCfg {
	return ServiceCfg{
		Port:        port,
		Env:         env,
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
	}
}

package config

type AppCfg struct {
	Port        string
	MaxChatSize int
	MaxChats    int
}

func AppInit(port string, maxChatSize int, maxChats int) AppCfg {
	return AppCfg{
		Port:        port,
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
	}
}

package proxy

import (
	"github.com/fuyao-w/sd/conf"
	"github.com/fuyao-w/sd/proxy/client"
	"github.com/fuyao-w/sd/proxy/server"
	"github.com/fuyao-w/sd/sd"
	"log"
)

var (
	DefaultConfig = Config{}
	SdComponent   sd.ServiceDiscover
)

type RegisterCenter struct {
	Type string `json:"type"`
	Addr string `json:"addr"`
}

type Config struct {
	Server         server.Server   `toml:"server"`
	Clients        []client.Client `toml:"clients"`
	ClientMap      map[string]client.Client
	RegisterCenter RegisterCenter `toml:"register_center"`
}

func Init() {
	conf.Decode("config/config.toml", &DefaultConfig)
	InitSd()
	InitClients()
}
func InitClients() {
	clientMap := make(map[string]client.Client, len(DefaultConfig.Clients))
	for _, client := range DefaultConfig.Clients {
		clientMap[client.Name] = client
	}
	DefaultConfig.ClientMap = clientMap
}

//func InitServer() *server.Server {
//	cfg := DefaultConfig.Server
//	s := &server.Server{
//		Name: cfg.Name,
//		Port: cfg.Port,
//	}
//	return s
//}

func InitSd() {
	var (
		cfg = DefaultConfig.RegisterCenter
		err error
	)
	switch cfg.Type {
	case "redis":
		SdComponent, err = sd.NewRedisRegisterProtocol(cfg.Addr)
	case "consul":
		log.Fatal("consul service register not implement")
		return
	}
	if err != nil {
		log.Fatalf("service init register err :%s\n", err.Error())
	}

}

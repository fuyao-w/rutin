package sd

import (
	"log"
)

var DefaultRegisterCenter ServiceDiscover

type RegisterCenter struct {
	Type string `json:"type"`
	Addr string `json:"addr"`
}

func InitSd(cfg RegisterCenter) {
	var (
		err error
	)
	switch cfg.Type {
	case "redis":
		DefaultRegisterCenter, err = NewRedisRegisterProtocol(cfg.Addr)
	case "consul":
		log.Fatal("consul service register not implement")
		return
	}
	if err != nil {
		log.Fatalf("service init register err :%s\n", err.Error())
	}
}

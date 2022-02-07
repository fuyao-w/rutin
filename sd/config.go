package sd

import (
	"context"
	"errors"
	"github.com/fuyao-w/sd/core"
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

func NewUpstream(serviceName string) core.Plugin {
	return core.Function(func(ctx context.Context, d core.Drive) {
		list := DefaultRegisterCenter.GetAddrSlice(serviceName)
		if list == nil {
			d.AbortErr(errors.New("no upstream"))
			return
		}

	})
}

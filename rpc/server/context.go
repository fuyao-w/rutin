package server

import (
	"context"
	"github.com/fuyao-w/rutin/core"
	"github.com/fuyao-w/rutin/discovery"
	"github.com/fuyao-w/rutin/load_balance"
	"github.com/fuyao-w/rutin/rpc/codec"
	"net"
)

type Plugin func(c *Context)

type Options struct {
	addr        net.Addr
	ServiceName string
	codec       codec.RequestCodec
	register    discovery.Register

	//handlerFunc iokit.HandlerFunc
}

type Option func(opt *Options)

func NewCodec(codec codec.RequestCodec) Option {
	return func(opt *Options) {
		opt.codec = codec
	}
}
func NewAddress(address string) Option {
	return func(opt *Options) {
		var err error
		opt.addr, err = net.ResolveTCPAddr("", address)
		if err != nil {
			panic("address invalid")
		}
	}
}
func WithLoadBalancer(t load_balance.LoadBalanceType) Option {
	return func(opt *Options) {
		load_balance.UpdateLocalLoadBalancer(t)
	}
}

func NewServiceName(name string) Option {
	return func(opt *Options) {
		opt.ServiceName = name
	}
}

type Context struct {
	driver     core.Driver
	opts       Options
	Ctx        context.Context
	Service    string
	Method     string
	RemoteAddr string
	Namespace  string
	Peer       string // 包含app_name的上游service_name
	Code       int32

	//// rpc request raw header
	//Header map[string]string

	// rpc request raw body
	Body     []byte
	Request  interface{}
	Response interface{}
}

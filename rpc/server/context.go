package server

import (
	"context"
	"github.com/fuyao-w/sd/core"
	"github.com/fuyao-w/sd/rpc/codec"
)

type Plugin func(c *Context)

type Options struct {
	Address     string
	ServiceName string
	codec       codec.RequestCodec
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
		opt.Address = address
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

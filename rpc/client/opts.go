package client

import (
	"github.com/fuyao-w/rutin/rpc/codec"
	"time"
)

type Options struct {
	//ServiceName string //简单通过服务发现名调用
	Codec      codec.RequestCodec
	PoolSize   int
	TTL        time.Duration
	RetryTimes int
	Timeout    time.Duration
	dialer     dialer
}

type Option func(options *Options)

func newOptions(opts ...Option) (options Options) {
	options = Options{
		//Codec:      nil,
		PoolSize:   DefaultPoolSize,
		TTL:        DefaultPoolTTL,
		RetryTimes: 1,
		Timeout:    3 * time.Second,
		Codec:      &codec.JsonCodec{},
	}
	options.dialer = &defaultDialer{options: options}
	for _, opt := range opts {
		opt(&options)
	}
	return
}

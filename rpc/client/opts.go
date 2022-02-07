package client

import "github.com/fuyao-w/sd/rpc/codec"

type Options struct {
	Codec      codec.RequestCodec
	PoolSize   int
	RetryTimes int
}

type Option func(options *Options)

func newOptions(opts ...Option) (options Options) {
	for _, opt := range opts {
		opt(&options)
	}
	return
}

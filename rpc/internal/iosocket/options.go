package iosocket

import "time"

type Options struct {
	RequestTimeout time.Duration
}

type Option func(opt *Options)

func Timeout(t time.Duration) Option {
	return func(opt *Options) {
		opt.RequestTimeout = t
	}
}

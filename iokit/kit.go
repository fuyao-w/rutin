package iokit

func NewCodec(codec MsgCodec) Option {
	return func(opt *options) {
		opt.codec = codec
	}
}

func NewOnMessage(f HandlerFunc) Option {
	return func(opt *options) {
		opt.handlerEntry = handlerEntry{HandlerFunc: f}
	}
}

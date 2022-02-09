package iosocket

import "sync"

type RequestContext struct {
	SeqID   uint64
	Request *Body
	Resp    *Body
	Closed  bool
	Err     error
	sync.WaitGroup
}

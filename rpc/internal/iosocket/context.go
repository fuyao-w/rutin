package iosocket

import (
	"github.com/fuyao-w/rutin/rpc/internal/metadata"
	"sync"
)

type RequestContextKey struct{}
type RequestContext struct {
	SeqID   uint64
	Request *metadata.HandlerDesc
	Resp    *Body
	Closed  bool
	Err     error
	sync.WaitGroup
}

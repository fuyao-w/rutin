package iosocket

import (
	"github.com/fuyao-w/rutin/rpc/internal/metadata"
	"sync"
	"time"
)

type RequestContextKey struct{}
type RequestContext struct {
	End     int64
	SeqID   uint64
	Request *metadata.HandlerDesc
	Resp    *Body
	Closed  bool
	Err     error
	sync.WaitGroup
	Timer *time.Timer
}

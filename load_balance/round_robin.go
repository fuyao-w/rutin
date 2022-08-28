package load_balance

import (
	"github.com/fuyao-w/rutin/discovery"
	"log"
	"sync/atomic"
)

type roundRobin struct {
	collection discovery.Collection
	index      uint64
}

func (r *roundRobin) Pick() discovery.Instance {
	newIdx := atomic.AddUint64(&r.index, 1)
	list := r.collection.GetInstances()
	if len(list) == 0 {
		log.Println("roundRobin.Pick list zero")
		return nil
	}
	return list[newIdx%uint64(len(list))]
}

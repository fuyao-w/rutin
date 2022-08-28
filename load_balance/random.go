package load_balance

import (
	"github.com/fuyao-w/rutin/discovery"
	"math/rand"
)

type randomPicker struct {
	collection discovery.Collection
}

func (r *randomPicker) Pick() discovery.Instance {
	list := r.collection.GetInstances()
	count := len(list)
	if count == 0 {
		return nil
	}
	return list[rand.Intn(count)]
}

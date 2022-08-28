package load_balance

import (
	"fmt"
	"github.com/fuyao-w/rutin/discovery"
)

const (
	LbRandom     LoadBalanceType = "random"
	LbRoundRobin LoadBalanceType = "round_robin"
)

type (
	LoadBalanceType string
	Picker          interface {
		Pick() discovery.Instance
	}
	pickerFunc func(discovery.Collection) Picker
)

var lbFuncMap = map[LoadBalanceType]pickerFunc{
	LbRandom: func(collection discovery.Collection) Picker {
		return &randomPicker{collection: collection}
	},
	LbRoundRobin: func(collection discovery.Collection) Picker {
		return &roundRobin{
			collection: collection,
		}
	},
}

type LoadBalancer interface {
	GetPicker(discovery.Collection) Picker
}

type LoadBalance struct {
	pickerFunc
}

func (l *LoadBalance) GetPicker(collection discovery.Collection) Picker {
	return l.pickerFunc(collection)
}

func NewLoadBalancer(loadBalanceType ...LoadBalanceType) LoadBalancer {
	if len(loadBalanceType) == 0 {
		loadBalanceType = append(loadBalanceType, LbRoundRobin)
	}
	f, ok := lbFuncMap[loadBalanceType[0]]
	if !ok {
		panic(fmt.Errorf("NewLoadBalancer not match load balance type :%s", string(loadBalanceType[0])))
	}
	return &LoadBalance{
		pickerFunc: f,
	}
}

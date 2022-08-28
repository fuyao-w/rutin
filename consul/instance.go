package consul

import (
	"github.com/fuyao-w/rutin/discovery"
	"net"
)

type ins struct {
	addr   net.Addr
	tags   []string
	weight int
}
type insCollection struct {
	list []discovery.Instance
}

func (i *insCollection) GetInstances() []discovery.Instance {
	return i.list
}

func (i *ins) GetAddr() net.Addr {
	return i.addr
}

func (i *ins) GetTags() []string {
	return i.tags
}

func (i *ins) Weight() int {
	return i.weight
}

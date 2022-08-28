package discovery

import (
	"net"
)

type RegisterCenter struct {
	Type string `json:"type"`
	Addr string `json:"addr"`
}

type Instance interface {
	GetAddr() net.Addr
	GetTags() []string
	Weight() int
}

type Collection interface {
	GetInstances() []Instance
}

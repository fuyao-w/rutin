package discovery

import (
	"github.com/fuyao-w/rutin/core"
	"github.com/fuyao-w/rutin/endpoint"
)

type ServiceDiscover interface {
	GetCollection(info *endpoint.ServiceInfo) Collection
}
type Register interface {
	Register(*endpoint.ServiceInfo) error
	Deregister(*endpoint.ServiceInfo) error
}

type PluginFactory interface {
	Factory(host string) (core.Plugin, error)
}

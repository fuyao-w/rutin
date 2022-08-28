package client

import (
	"github.com/fuyao-w/rutin/consul"
	"github.com/fuyao-w/rutin/endpoint"
)

//type RClient struct {
//	Name          string `json:"name"`
//	//EndpointsFrom string `json:"endpoints_from"` //redis consul
//}
//最外层初始化的时候用
func RpcClient(endpoint *endpoint.ServiceInfo, options ...Option) NetClient {
	opts := newOptions(options...)
	opts.discovery = consul.NewConsulDiscovery()
	return newGeneralClient(&rpcFactory{
		options:  opts,
		connPool: newPool(opts.PoolSize, opts.TTL, opts.dialer),
	}, endpoint, opts)
}

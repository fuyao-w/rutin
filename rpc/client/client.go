package client

import (
	"context"
	"github.com/fuyao-w/rutin/core"
	"github.com/fuyao-w/rutin/discovery"
	"github.com/fuyao-w/rutin/endpoint"
	"github.com/fuyao-w/rutin/kit/recovery"
	"github.com/fuyao-w/rutin/kit/retry"
	"github.com/fuyao-w/rutin/kit/upstream"
	"log"
)

type NetClient interface {
	Invoke(ctx context.Context, methName string, in, out interface{}) error
}

type generalClient struct {
	options       Options
	serviceName   string //调用的 RPC 方法
	defaultDriver core.Drive
}

func (g *generalClient) Invoke(ctx context.Context, methName string, in, out interface{}) error {
	driver := g.defaultDriver.Copy()
	driver.Next(context.WithValue(ctx, rpcContextKey, &RpcContext{
		EndPoint: methName,
		Request:  in,
		Response: out,
	}))
	return driver.Err()
}

func newGeneralClient(factory discovery.PluginFactory, endpoint *endpoint.ServiceInfo, options Options) *generalClient {
	plugins := []core.Plugin{
		recovery.Recover(),
		retry.Retry(options.RetryTimes),
		upstream.NewUpStream(factory, func() discovery.Collection {
			list := options.discovery.GetCollection(endpoint)
			if len(list.GetInstances()) == 0 {
				log.Println("len(list.GetInstances()) == 0 ")
			}
			return list
		}),
	}

	return &generalClient{
		options:       options,
		serviceName:   endpoint.Name,
		defaultDriver: core.New(plugins),
	}
}

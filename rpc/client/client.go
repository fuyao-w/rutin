package client

import (
	"context"
	"github.com/fuyao-w/sd/core"
	"github.com/fuyao-w/sd/sd"
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
		ServiceName: g.serviceName,
		EndPoint:    methName,
		Request:     in,
		Response:    out,
	}))
	return driver.Err()
}

func newGeneralClient(factory sd.PluginFactory, serviceName string, options Options) *generalClient {
	plugins := []core.Plugin{
		//recovery.Recover(),
		//retry.Retry(options.RetryTimes),
		sd.NewUpStream(factory, serviceName),
	}

	return &generalClient{
		options:       options,
		serviceName:   serviceName,
		defaultDriver: core.New(plugins),
	}
}

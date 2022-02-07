package client

import (
	"context"
	"github.com/fuyao-w/sd/core"
	"github.com/fuyao-w/sd/kit/recovery"
	"github.com/fuyao-w/sd/kit/retry"
)

type NetClient interface {
	Invoke(ctx context.Context, methName string, in, out interface{}) error
}

type generalClient struct {
	options       Options
	RpcMeth       string //调用的 RPC 方法
	defaultDriver core.Drive
}

func NewGeneralClient(rpcMeth string, options Options) *generalClient {
	plugins := []core.Plugin{
		recovery.Recover(),
		retry.Retry(options.RetryTimes),
	}

	return &generalClient{
		options:       options,
		RpcMeth:       rpcMeth,
		defaultDriver: core.New(plugins),
	}
}

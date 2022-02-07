package client

import (
	"context"
	"github.com/fuyao-w/sd/core"
)

type rpcFactory struct {
	options  Options
	connPool pool
}

func (r *rpcFactory) Factory(host string) (core.Plugin, error) {
	var (
		codec = r.options.Codec
	)
	socket, err := r.connPool.GetSocket(host)
	if err != nil {
		return nil, err
	}
	return core.Function(func(ctx context.Context, core core.Drive) {
		rpcCtx := ctx.Value(rpcContextKey).(RpcContext)
		body, err := codec.Encode(rpcCtx.Request)
		if err != nil {
			core.AbortErr(err)
		}

		body, err = socket.Call(rpcCtx.EndPoint, body)
		if err != nil {

		}
	}), nil
}

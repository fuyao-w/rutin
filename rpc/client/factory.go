package client

import (
	"context"
	"github.com/fuyao-w/rutin/core"
	"github.com/fuyao-w/rutin/rpc/internal/iosocket"
	"log"
)

type rpcFactory struct {
	options  Options
	connPool pool
}

func (r *rpcFactory) Factory(host string) (core.Plugin, error) {
	socket, err := r.connPool.getSocket(host)
	if err != nil {
		log.Printf("rpcFactory|Factory|getSocket err %s,host: %s", err, host)
		return nil, err
	}
	return core.Function(func(ctx context.Context, core core.Drive) {
		var (
			err  error
			body []byte
		)
		defer func() {
			if err != nil {
				core.AbortErr(err)
			}
		}()
		rpcCtx := ctx.Value(rpcContextKey).(*RpcContext)

		body, err = socket.Call(rpcCtx.EndPoint, rpcCtx.Request)
		switch {
		case err == iosocket.InternalErrExited:
			r.connPool.releaseSock(host, socket)
			return
		case err != nil:
			return
		}
		err = r.options.Codec.Decode(body, rpcCtx.Response)
	}), err
}

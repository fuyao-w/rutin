package client

import (
	"encoding/json"
	"github.com/fuyao-w/sd/rpc/codec"
	"github.com/fuyao-w/sd/rpc/internal/metadata"
)

type socket interface {
	Call(endpoint /*rpc 的调用方法*/ string, body []byte) ([]byte, error)
	Close() error
}

type RpcSocket struct {
	codec codec.RequestCodec
}

func (r *RpcSocket) Call(endpoint string, body []byte) ([]byte, error) {
	metaDate := metadata.HandlerDesc{
		MethName: endpoint,
		Param:    body,
	}


}

func (r *RpcSocket) Close() error {
	panic("implement me")
}

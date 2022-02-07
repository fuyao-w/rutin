package client

import (
	"context"
	"encoding/json"
	"github.com/fuyao-w/sd/io/codec"
	"github.com/fuyao-w/sd/rpc/internal/iosocket"
	"github.com/fuyao-w/sd/rpc/internal/metadata"

	"github.com/fuyao-w/sd/utils"
	"log"
)

type Client struct {
	Name          string `json:"name"`
	EndpointsFrom string `json:"endpoints_from"` //redis consul
}

func (c *Client) Call(ctx context.Context, msg []byte) ([]byte, error) {
	rpcCtx := ctx.Value(rpcContextKey).(*RpcContext)

	return iosocket.Client(rpcCtx.EndPoint, msg, &codec.ProtocolParser{})
}

func (c *Client) Invoke(ctx context.Context, methName string, req, reply interface{}) error {

	desc := metadata.HandlerDesc{
		MethName: methName,
		Param:    utils.GetJsonBytes(req),
	}

	log.Println("producer body", string(utils.GetJsonBytes(desc)), string(utils.GetJsonBytes(req)))
	resp, err := c.Call(ctx, utils.GetJsonBytes(desc))
	if err != nil {
		log.Println("Producer err", err)
		return err
	}
	err = json.Unmarshal(resp, &reply)
	if err != nil {
		log.Println("proxy calc err", err)
	}
	return err
}

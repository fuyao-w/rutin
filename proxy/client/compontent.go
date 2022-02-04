package client

import (
	"errors"
	"github.com/fuyao-w/sd/net"
	"github.com/fuyao-w/sd/parse"
	"github.com/fuyao-w/sd/proxy"
	"log"
)

type Client struct {
	Name          string `json:"name"`
	EndpointsFrom string `json:"endpoints_from"` //redis consul
}



func (c *Client) Call(msg []byte) ([]byte, error) {
	addrs := proxy.SdComponent.GetAddrSlice(c.Name)
	if len(addrs) == 0 {
		log.Println("no upstream")
		return nil, errors.New("no upstream")
	}
	return net.Client(addrs[0], msg, &parse.ProtocolParser{})
}
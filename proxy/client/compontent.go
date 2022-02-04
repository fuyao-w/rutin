package client

import (
	"errors"
	"github.com/fuyao-w/sd/net"
	"github.com/fuyao-w/sd/parse"
	"github.com/fuyao-w/sd/sd"
	"log"
)

type Client struct {
	Name          string `json:"name"`
	EndpointsFrom string `json:"endpoints_from"` //redis consul
	sdComponent   sd.ServiceDiscover
}

func (c *Client) Init() {
	var err error
	switch c.EndpointsFrom {
	case "redis":
		c.sdComponent, err = sd.NewRedisRegisterProtocol()
	case "consul":
		log.Fatal("consul service discover not implement")
	}

	if err != nil {
		log.Println("init err ", err)
		return
	}
	c.sdComponent.SetRegisterName(c.Name)
}

func (c *Client) Call(msg []byte) ([]byte, error) {
	addrs := c.sdComponent.GetAddrSlice()
	if len(addrs) == 0 {
		log.Println("no upstream")
		return nil, errors.New("no upstream")
	}
	return net.Client(addrs[0], msg, &parse.ProtocolParser{})
}

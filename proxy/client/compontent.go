package client

import (
	"encoding/json"
	"errors"
	"github.com/fuyao-w/sd/net"
	"github.com/fuyao-w/sd/parse"
	"github.com/fuyao-w/sd/proxy/server"
	"github.com/fuyao-w/sd/sd"
	"github.com/fuyao-w/sd/utils"
	"log"
)

type Client struct {
	Name          string `json:"name"`
	EndpointsFrom string `json:"endpoints_from"` //redis consul
}

func (c *Client) Call(msg []byte) ([]byte, error) {
	addrs := sd.DefaultRegisterCenter.GetAddrSlice(c.Name)
	if len(addrs) == 0 {
		log.Println("no upstream")
		return nil, errors.New("no upstream")
	}
	return net.Client(addrs[0], msg, &parse.ProtocolParser{})
}

func Register(serviceName, methName string, req, reply interface{}) error {
	client := DefaultClientMap[serviceName]
	desc := server.HandlerDesc{
		ServiceName: serviceName,
		MethName:    methName,
		Param:       utils.GetJsonBytes(req),
	}

	log.Println("producer body", string(utils.GetJsonBytes(desc)), string(utils.GetJsonBytes(req)))
	resp, err := client.Call(utils.GetJsonBytes(desc))
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

package worker

import (
	"encoding/json"
	"github.com/fuyao-w/sd/proxy/client"
	"github.com/fuyao-w/sd/proxy/server"
	"github.com/fuyao-w/sd/utils"
	"log"
)

type ProxyHandle struct {
	client *client.Client
}

func InitProxyHandle() *ProxyHandle {
	c := &client.Client{
		Name:          "calc",
		EndpointsFrom: "redis",
	}
	c.Init()
	return &ProxyHandle{
		client: c,
	}
}

func InitServer() *server.Server {
	handle := Handle{}
	s := &server.Server{
		Name: handle.Name(),
		Port: 10010,
		RegisterCenter: server.RegisterCenter{
			Type: "redis",
			Addr: "127.0.0.1:6379",
		},
	}
	server.RegisterHandle(&handle)
	s.Init()
	return s
}

type Handler interface {
	Calc(req ClacReq, resp *ClacResp) error
}
type ClacReq struct {
	A int `json:"a"`
	B int `json:"b"`
}
type ClacResp struct {
	DmError int `json:"dm_error"`
	Result  int `json:"result"`
}

func (p *ProxyHandle) Calc(req ClacReq, calcResp *ClacResp) error {
	desc := server.HandlerDesc{
		ServiceName: p.client.Name,
		MethName:    "Calc",
		Param:       utils.GetJsonBytes(req),
	}

	log.Println("producer body", string(utils.GetJsonBytes(desc)), string(utils.GetJsonBytes(req)))
	resp, err := p.client.Call(utils.GetJsonBytes(desc))
	if err != nil {
		log.Println("Producer err", err)
		calcResp = &ClacResp{
			DmError: 500,
		}
		return nil
	}
	err = json.Unmarshal(resp, &calcResp)
	if err != nil {
		log.Println("proxy calc err", err)
	}

	return nil
}

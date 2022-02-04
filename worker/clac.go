package worker

import (
	"encoding/json"
	"github.com/fuyao-w/sd/proxy"

	"github.com/fuyao-w/sd/proxy/server"
	"github.com/fuyao-w/sd/utils"
	"log"
)

type ProxyHandle struct {
	//client *client.Client
}

func (p *ProxyHandle) Name() string {
	return "calc"
}

func InitProxyHandle(name string) *ProxyHandle {
	return &ProxyHandle{}
}

func InitHandle() (handle *Handle) {
	handle = &Handle{}
	server.RegisterHandle(handle)
	return handle
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
	client := proxy.DefaultConfig.ClientMap[p.Name()]
	desc := server.HandlerDesc{
		ServiceName: p.Name(),
		MethName:    "Calc",
		Param:       utils.GetJsonBytes(req),
	}

	log.Println("producer body", string(utils.GetJsonBytes(desc)), string(utils.GetJsonBytes(req)))
	resp, err := client.Call(utils.GetJsonBytes(desc))
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

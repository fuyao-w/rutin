package worker

import (
	"github.com/fuyao-w/sd/rpc/proxy/client"
	"github.com/fuyao-w/sd/rpc/proxy/server"
)

type ProxyHandle struct {
	//client *client.Client
}

const serviceName = "calc"

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
	return client.Call(p.Name(), "Calc", req, calcResp)
}

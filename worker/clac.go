package worker

import (
	"context"
	"github.com/fuyao-w/rutin/endpoint"
	"github.com/fuyao-w/rutin/rpc/client"
)

type ProxyHandle struct {
	client client.NetClient
}

const serviceName = "user.account.login"

func (p *ProxyHandle) Name() string {
	return "calc"
}

func InitProxyHandle() *ProxyHandle {
	return &ProxyHandle{
		client: client.RpcClient(&endpoint.ServiceInfo{
			Name: serviceName,
		}),
	}
}

func InitHandle() (handle *Handle) {
	handle = &Handle{}
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
	return p.client.Invoke(context.TODO(), serviceName+"|"+"Calc", req, calcResp)
}

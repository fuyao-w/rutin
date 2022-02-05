package main

import (
	"fmt"
	"github.com/fuyao-w/sd/proxy/client"
	"github.com/fuyao-w/sd/proxy/server"
	"github.com/fuyao-w/sd/sd"
	"github.com/fuyao-w/sd/worker"
)

func Init() {
	sd.InitSd(sd.RegisterCenter{
		Type: "redis",
		Addr: "127.0.0.1:6379",
	})
	client.NewClients(client.Client{
		Name:          "calc",
		EndpointsFrom: "redis",
	})
	server.ConfigServer(server.Server{
		Name: "calc",
		Port: 10010,
	})
}
func main() {
	var (
		c = worker.InitProxyHandle("calc")
	)
	Init()
	worker.InitHandle() //server 注册
	go server.BeginServer()
	reply := &worker.ClacResp{}
	err := c.Calc(worker.ClacReq{
		A: 1,
		B: 2,
	}, reply)
	fmt.Println(err, reply.DmError, reply.Result)

}

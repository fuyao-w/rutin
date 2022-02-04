package main

import (
	"fmt"
	"github.com/fuyao-w/sd/proxy"

	"github.com/fuyao-w/sd/proxy/server"
	"github.com/fuyao-w/sd/worker"
)

func main() {
	var (
		client = worker.InitProxyHandle("calc")
	)
	proxy.Init()
	worker.InitHandle()
	go server.BeginServer()
	reply := &worker.ClacResp{}
	err := client.Calc(worker.ClacReq{
		A: 1,
		B: 2,
	}, reply)
	fmt.Println(err, reply.DmError, reply.Result)

}

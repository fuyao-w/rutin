package main

import (
	"fmt"
	"github.com/fuyao-w/sd/worker"
)

func main() {
	var (
		client = worker.InitProxyHandle()
		server = worker.InitServer()
	)
	go server.Server()
	reply := &worker.ClacResp{}
	err := client.Calc(worker.ClacReq{
		A: 1,
		B: 2,
	}, reply)
	fmt.Println(err, reply.DmError, reply.Result)

}

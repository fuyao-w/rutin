package main

import (
	"fmt"
	"github.com/fuyao-w/rutin/rpc/codec"
	"github.com/fuyao-w/rutin/rpc/server"
	"github.com/fuyao-w/rutin/sd"
	"github.com/fuyao-w/rutin/worker"
	"time"
)

func init() {
	sd.InitSd(sd.RegisterCenter{
		Type: "redis",
		Addr: "127.0.0.1:6379",
	})
	//client.NewClients(client.Client{
	//	Name:          "calc",
	//	EndpointsFrom: "redis",
	//})
	//server.ConfigServer(server.Server{
	//	Name: "calc",
	//	Port: 10010,
	//})
}
func main() {
	//client := worker.InitProxyHandle()
	go func() {
		handler := worker.InitHandle()
		server := server.NewRpcServer(
			server.NewAddress("127.0.0.1:10000"),
			server.NewCodec(&codec.JsonCodec{}),
		)
		server.Handle(server.NewHandler(handler))
		fmt.Println(server.GetPaths())
		server.Start()
	}()
	time.Sleep(time.Millisecond * 40)
	handle := worker.InitProxyHandle()
	var resp worker.ClacResp
	for i := 0; i < 100; i++ {

		fmt.Println(handle.Calc(worker.ClacReq{
			A: 1,
			B: 2,
		}, &resp))

		fmt.Printf("%d resp :%+v\n", i, resp)
	}

	//go func() {
	//
	//}()

}

package main

import (
	"fmt"
	"github.com/fuyao-w/rutin/rpc/codec"
	"github.com/fuyao-w/rutin/rpc/server"
	"github.com/fuyao-w/rutin/sd"
	"github.com/fuyao-w/rutin/worker"
	"log"
	"sync"
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
	time.Sleep(time.Millisecond * 100)
	handle := worker.InitProxyHandle()
	var (
		wg  sync.WaitGroup
		sum = 2200
	)
	now := time.Now().UnixNano()
	wg.Add(sum)
	for i := 0; i < sum; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			var resp worker.ClacResp
			//log.Printf("%d resp :%+v\n", i, resp)
			if err := handle.Calc(worker.ClacReq{
				A: i,
				B: i >> 1,
			}, &resp); err != nil {
				log.Printf("calc err %s", err)
			} else {
				//fmt.Println("resp", resp.Result)
			}


		}(i)

	}
	wg.Wait()
	log.Printf("execute time :%f\n\n\n", time.Duration(time.Now().UnixNano()-now).Seconds())
	//go func() {
	//
	//}()

}

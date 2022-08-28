package main

import (
	"fmt"
	"github.com/fuyao-w/rutin/rpc/codec"
	"github.com/fuyao-w/rutin/rpc/server"
	"github.com/fuyao-w/rutin/worker"
	"log"
	"sync"
	"time"
)

const serviceName = "user.account.login"

func main() {
	//client := worker.InitProxyHandle()
	go func() {
		name := server.WithServiceName(serviceName)
		handler := worker.InitHandle()
		server := server.NewRpcServer(
			server.NewAddress("127.0.0.1:10000"),
			server.NewCodec(&codec.JsonCodec{}),
		)
		server.Handle(server.NewHandler(handler, name))
		fmt.Println(server.GetPaths())
		server.Start()
	}()
	time.Sleep(time.Second * 1)
	handle := worker.InitProxyHandle()
	var (
		wg  sync.WaitGroup
		sum = 100
	)
	now := time.Now().UnixNano()
	wg.Add(sum)
	for i := 0; i < sum; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			var resp worker.ClacResp
			//log.Printf("%d r esp :%+v\n", i, resp)
			if err := handle.Calc(worker.ClacReq{
				A: i,
				B: i >> 1,
			}, &resp); err != nil {
				log.Printf("calc err %s", err)
			} else {
				fmt.Println("resp", resp.Result)
			}

		}(i)

	}
	wg.Wait()
	time.Sleep(time.Hour)
	log.Printf("execute time :%f\n\n\n", time.Duration(time.Now().UnixNano()-now).Seconds())
	//go func() {
	//
	//}()

}

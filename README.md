# rutin
里面目前有简陋的 rpc  、通过服务发现调用

目前没有自动生成静态代理 handler 的功能，只能自己手写

示例：

```go

func init() {
	sd.InitSd(sd.RegisterCenter{
		Type: "redis",
		Addr: "127.0.0.1:6379",
	})

}
func main() {
	
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
}

```

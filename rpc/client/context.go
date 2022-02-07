package client

var rpcContextKey = "rpc_context_key"

type RpcContext struct {
	EndPoint          string
	Request, Response interface{}
	retryMax          int
	retry             bool
	//host              string
}

//func (receiver ) name()  {
//
//}
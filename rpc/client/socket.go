package client

import (
	"errors"
	"github.com/fuyao-w/rutin/rpc/codec"
	"github.com/fuyao-w/rutin/rpc/internal/iosocket"
	"github.com/fuyao-w/rutin/rpc/internal/metadata"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

type socket interface {
	Call(endpoint /*rpc 的调用方法*/ string, req interface{}) ([]byte, error)
	Close() error
}

type RpcSocket struct {
	timeOut time.Duration
	socket  iosocket.IoSocket
	codec   codec.RequestCodec
}

func initRpcSocket(conn net.Conn, options Options) *RpcSocket {
	s := &RpcSocket{
		timeOut: options.Timeout,
		codec:   options.Codec,
	}
	socket := iosocket.NewIoSocket(conn, options.Codec)
	socket.StartWorker()
	s.socket = socket

	return s
}

func (r *RpcSocket) Call(endpoint string, req interface{}) ([]byte, error) {
	arr := strings.Split(endpoint, ".")
	if len(arr) != 2 {
		return nil, errors.New("endpoint format err .right format : serviceName.methName")
	}

	payload, err := r.codec.Encode(req)
	if err != nil {
		log.Printf("RpcSocket|Call|Encode err %s", err)
		return nil, err
	}
	body, err := r.socket.Call(metadata.HandlerDesc{
		ServiceName: arr[0],
		MethName:    arr[1],
		Param:       payload,
	}, atomic.AddUint64(&r.socket.SeqID, 1))

	if err != nil {
		log.Printf("RpcSocket|Call err %s", err)
		return nil, err
	}
	//log.Printf("RpcSocket|Call Payload %s", body.Payload)
	return body.Payload, err
}

func (r *RpcSocket) Close() error {
	return r.socket.Close()
}

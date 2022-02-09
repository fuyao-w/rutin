package client

import (
	"github.com/fuyao-w/sd/rpc/codec"
	"github.com/fuyao-w/sd/rpc/internal/iosocket"
	"github.com/fuyao-w/sd/rpc/internal/metadata"
	"log"
	"strings"

	"net"
	"time"
)

type socket interface {
	Call(endpoint /*rpc 的调用方法*/ string, body []byte) ([]byte, error)
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

func (r *RpcSocket) Call(endpoint string, payload []byte) ([]byte, error) {
	arr := strings.Split(endpoint, ".")
	metaDate := metadata.HandlerDesc{
		ServiceName: arr[0],
		MethName:    arr[1],
		Param:       payload,
		//SeqID:       r.socket.SeqID,
	}
	bytes, err := r.codec.Encode(metaDate)
	if err != nil {
		return nil, err
	}

	body, err := r.socket.Call(&iosocket.Body{Payload: bytes})
	if err != nil {

		return nil, err
	}
	log.Printf("RpcSocket|Call err %s", body.Payload)
	return body.Payload, err
}

func (r *RpcSocket) Close() error {
	return r.socket.Close()
}

package iosocket

import (
	"errors"
	"github.com/fuyao-w/sd/iokit"
	"github.com/fuyao-w/sd/rpc/codec"
	"github.com/fuyao-w/sd/rpc/internal/metadata"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

var (
	ErrExited   = errors.New("socket exit")
	ErrChanSize = errors.New("socket chan blocked")
)

type IoSocket struct {
	closed     int64
	exitC      chan struct{}
	SeqID      uint64
	conn       net.Conn
	peerAddr   string
	localAddr  string
	requestC   chan *RequestContext
	ExistC     chan struct{}
	wc         io.WriteCloser
	codec      codec.RequestCodec
	controller sync.Map
	serverCB   iokit.HandlerFunc
}

type Body struct {
	Payload []byte
}

func NewIoSocket(conn net.Conn, codec codec.RequestCodec) IoSocket {
	return IoSocket{
		codec:    codec,
		conn:     conn,
		requestC: make(chan *RequestContext, 1024),
		ExistC:   make(chan struct{}),
	}
}

//func (s *IoSocket) Start(host string, timeOut time.Duration) (conn net.Conn, err error) {
//	conn, err = net.DialTimeout("tcp", host, timeOut)
//	if err != nil {
//		log.Printf("IoSocket|Start|DialTimeout err :%s,host :%s,", err, host)
//		return nil, err
//	}
//	return
//}

func (s *IoSocket) StartWorker() {

	client := iokit.NewClientChannel(s.conn, iokit.NewOnMessage(s.ClientOnMessage), iokit.NewCodec(&iokit.ProtocolParser{}))
	client.Start()
	s.wc = client
	go s.dispatch()
}
func (s *IoSocket) dispatch() {
	for {
		select {
		case <-s.ExistC:
			return
		case req := <-s.requestC:
			s.controller.Store(req.SeqID, req)
			pckParser := iokit.SeqPacket{SeqID: req.SeqID, Payload: req.Request.Payload}
			body, err := pckParser.Encode()
			if err != nil {
				return
			}
			s.wc.Write(body)
		}
	}
}

func (s *IoSocket) Call(handlerDesc metadata.HandlerDesc) (*Body, error) {
	bytes, err := s.codec.Encode(handlerDesc)
	if err != nil {
		log.Printf("RpcSocket|Call|Encode err %s", err)
		return nil, err
	}
	request := &RequestContext{
		SeqID:   handlerDesc.SeqID,
		Request: &Body{Payload: bytes},
	}
	request.Add(1)
	select {
	case <-s.exitC:
		return nil, ErrExited
	case s.requestC <- request:
	default:
		return nil, ErrChanSize
	}
	request.Wait()

	return request.Resp, nil
}
func (i *IoSocket) Close() error {
	if !atomic.CompareAndSwapInt64(&i.closed, 0, 1) {
		// already close
		return nil
	}

	close(i.exitC)

	// we will not close requestC chan here.
	// close(i.requestC)

	//i.wg.Wait()

	//还要先 cancel 所有的请求

	//i.controller.Range(func(key, value interface{}) bool {
	//	request := value.(*requestContext)
	//	i.contextEnd(request, nil, ErrExited)
	//	return true
	//})
	return i.wc.Close()
}

func (i *IoSocket) ClientOnMessage(body []byte, _ io.WriteCloser) {
	seq := iokit.SeqPacket{}
	if err := seq.Decode(body); err != nil {
		log.Printf("IoSocket|ClientOnMessage err %s", err)
		return
	}
	desc, err := metadata.Parse(i.codec, seq.Payload)
	if err != nil {
		log.Printf("IoSocket|ClientOnMessage metadata.Parse err:%s ,body :%s", err, string(body))
		return
	}

	value, ok := i.controller.Load(desc.SeqID)
	if !ok {
		log.Printf("IoSocket|ClientOnMessage|Load not ok :%+v", desc)
		return
	}

	ctx, _ := value.(*RequestContext)

	i.controller.Delete(seq.SeqID)
	i.requestCtxEnd(ctx, desc.Response)

}

func (i *IoSocket) requestCtxEnd(ctx *RequestContext, body []byte) {
	ctx.Resp = &Body{Payload: body}
	ctx.Request = nil
	ctx.Done()
}

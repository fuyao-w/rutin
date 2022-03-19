package iosocket

import (
	"errors"
	"github.com/fuyao-w/rutin/iokit"
	"github.com/fuyao-w/rutin/rpc/codec"
	"github.com/fuyao-w/rutin/rpc/internal/metadata"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	InternalErrExited   = errors.New("iosocket: socket exit")
	InternalErrChanSize = errors.New("iosocket: socket chan blocked")
	ErrRequestTimeOut   = errors.New("iosocket: request time out")
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
			//log.Printf("dispatch id :%d", req.SeqID)
			s.controller.Store(req.SeqID, req)
			body, _ := (&iokit.SeqPacket{
				SeqID: req.SeqID,
				Payload: func() []byte {
					body, _ := metadata.Marshal(req.Request)
					return body
				}(),
			}).Encode()

			//if len(body) > 156 {
			//	fmt.Println("dispatch", len(body), string(body))
			//	parser := iokit.ProtocolParser{}
			//	d,_ := parser.Encode(body)
			//	fmt.Println(string(d))
			//}

			if _, err := s.wc.Write(body); err != nil {
				s.requestCtxEnd(req, nil, err)
			}
		}
	}
}

/*

 */
func (s *IoSocket) Call(handlerDesc metadata.HandlerDesc, seqID uint64, opts ...Option) (*Body, error) {
	var options = Options{
		RequestTimeout: time.Second,
	}
	for _, opt := range opts {
		opt(&options)
	}

	request := &RequestContext{
		SeqID:   seqID,
		Request: &handlerDesc,
	}

	request.Add(1)
	request.Timer = time.AfterFunc(options.RequestTimeout, func() {
		s.requestCtxEnd(request, nil, ErrRequestTimeOut)
	})
	select {
	case <-s.exitC: //socket 关闭
		return nil, InternalErrExited
	case s.requestC <- request: //正常完成

	default: //请求过多
		return nil, InternalErrChanSize
	}

	request.Wait()

	return request.Resp, request.Err
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
	//	i.contextEnd(request, nil, InternalErrExited)
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
	desc, err := metadata.Unmarshal(seq.Payload)
	if err != nil {
		log.Printf("IoSocket|ClientOnMessage metadata.Unmarshal err:%s ,body :%s", err, string(body))
		return
	}

	value, ok := i.controller.Load(seq.SeqID)
	if !ok {
		log.Printf("IoSocket|ClientOnMessage|Load not ok :%+v", desc)
		return
	}

	ctx, _ := value.(*RequestContext)
	//fmt.Println("values", *ctx)
	i.controller.Delete(ctx.SeqID)

	i.requestCtxEnd(ctx, desc.Response, nil)

}

func (i *IoSocket) requestCtxEnd(ctx *RequestContext, body []byte, err error) {
	//这块防止 requestCtxEnd 重复调用导致 WaitGroup.Done() panic
	if !atomic.CompareAndSwapInt64(&ctx.End, 0, 1) {
		// context already end
		return
	}
	ctx.Resp = &Body{Payload: body}
	ctx.Request = nil
	ctx.Err = err
	ctx.Timer.Stop()
	ctx.Done()
}

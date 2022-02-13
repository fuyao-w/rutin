package iokit

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type (
	options struct {
		codec             MsgCodec
		handlerEntry      handlerEntry
		timerBufferSize   int // size of buffered channel
		handlerBufferSize int // size of buffered channel
		writerBufferSize  int // size of buffered channel
	}

	IoServer struct {
		options   options
		ctx       context.Context
		cancel    context.CancelFunc
		conns     sync.Map
		codec     MsgCodec
		listeners map[net.Listener]bool
		//messageRegistry map[int32]*handlerEntry -- 这个后面用
		messageRegistry *handlerEntry
		mu              sync.Mutex
		connID          uint64
	}

	HandlerFunc func(body []byte, wc io.WriteCloser)

	handlerEntry struct {
		HandlerFunc HandlerFunc
	}
	MsgHandler struct {
		Msg     []byte
		Handler handlerEntry
	}
	Option func(opt *options)
)

func NewIoServer(opts ...Option) *IoServer {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	return &IoServer{
		options:   options,
		ctx:       nil,
		cancel:    nil,
		conns:     sync.Map{},
		codec:     options.codec,
		listeners: map[net.Listener]bool{},
		//messageRegistry: map[int32]*handlerEntry{},
		mu:     sync.Mutex{},
		connID: 0,
	}
}
func (s *IoServer) Start(l net.Listener) error {
	var i = 1
	for {
		conn, err := l.Accept()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Temporary() {
				log.Printf("IoServer|Accept|err %s ,idx :%d", err, i)
				//time.Sleep(time.Second)
				continue
			}
			return err
		}
		i++
		channel := NewServerChannel(conn, atomic.AddUint64(&s.connID, 1), channelOptions{options: s.options})
		go channel.Start()

	}

}
func (s *IoServer) Stop() {
	for listener := range s.listeners {
		listener.Close()
	}
	s.cancel()
}

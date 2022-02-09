package iokit

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed = errors.New("server closed")
	ErrWouldBlock   = errors.New("writer chan blocked")
)

type channelOptions struct {
	options
}
type Channel struct {
	options   channelOptions
	cancelCtx context.Context
	cancel    context.CancelFunc
	conn      net.Conn
	connID    uint64
	reader    *bufio.Reader
	writer    *bufio.Writer
	//codec      MsgCodec
	lastActive int64
	closed     int32
	writeC     chan []byte
	handleC    chan *MsgHandler
	wg         sync.WaitGroup
	once       sync.Once //Close 的时候用
}

type ClientChannel struct {
	*Channel
}
type ServerChannel struct {
	*Channel
}

// Close asyncclose
func (c *ClientChannel) Close() error {
	c.cancel()
	return nil
}

// Close asyncclose
func (c *ServerChannel) Close() error {
	c.cancel()
	return nil
}

func (c *Channel) Write(p []byte) (n int, err error) {
	if atomic.LoadInt32(&c.closed) == 1 {
		return 0, ErrServerClosed
	}

	defer func() {
		if ppp := recover(); ppp != nil {
			err = ErrServerClosed
		}
	}()

	bytes, err := c.options.codec.Encode(p)
	if err != nil {
		return 0, err
	}

	select {
	case c.writeC <- bytes:
		err = nil
		//fmt.Println("writeloop" ,string(bytes))
	default:
		err = ErrWouldBlock
	}
	return len(p), err
}

func (c *Channel) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}
	c.once.Do(func() {
		fmt.Println("close")
		log.Fatal("channel clost")
		c.conn.Close()
		close(c.writeC)
		close(c.handleC)
	})
	return nil
}

func NewServerChannel(conn net.Conn, connID uint64, options channelOptions) *ServerChannel {

	return &ServerChannel{
		Channel: (&Channel{
			options: options,
			connID:  connID,
			conn:    conn,
			writeC:  make(chan []byte),
			handleC: make(chan *MsgHandler),
		}).Init(),
	}
}
func NewClientChannel(conn net.Conn, opts ...Option) *ClientChannel {
	var options options
	for _, o := range opts {
		o(&options)
	}
	return &ClientChannel{
		Channel: (&Channel{
			options: channelOptions{options: options},
			conn:    conn,
			writeC:  make(chan []byte),
			handleC: make(chan *MsgHandler),
		}).Init(),
	}
}

func (c *Channel) deferProc(name string, wg *sync.WaitGroup) {
	wg.Done()

	if p := recover(); p != nil {
		//c.Close()
		log.Printf("%s panic :%s", name, p)
		debug.PrintStack()
		return
	}
}
func (c *Channel) readLoop() {
	defer c.deferProc("readLoop", &c.wg)
	for {
		select {
		case <-c.cancelCtx.Done():
			return
		default:

			body, err := c.options.codec.Decode(c.conn)
			if err != nil {
				log.Printf("readLoop|Decode err %s", err)
				if err, ok := err.(net.Error); ok && err.Temporary() {
					time.Sleep(time.Second)
					continue
				}
				return

			}
			fmt.Printf("readloop :%s\n", c.options.handlerEntry.HandlerFunc == nil)

			c.handleC <- &MsgHandler{
				Msg:     body,
				Handler: c.options.handlerEntry,
			}
		}
	}
}
func (c *Channel) writeLoop() {
	defer c.deferProc("writeLoop", &c.wg)
	for {
		select {
		case <-c.cancelCtx.Done():
			return
		case info := <-c.writeC:
			//fmt.Println("resdfsdf\n", string(info),"\n")
			//body, err := c.options.codec.Encode(info)
			//if err != nil {
			//	log.Printf("handleLoop|Encode err %s", err)
			//	continue
			//}

			c.writer.Write(info)
			c.writer.Flush()

		}
	}
}
func (c *Channel) handleLoop() {
	defer c.deferProc("handleLoop", &c.wg)
	for {
		select {
		case <-c.cancelCtx.Done():
			return
		case m := <-c.handleC:

			m.Handler.HandlerFunc(m.Msg, c)
		}
	}
}

func (c *Channel) Start() {
	c.wg.Add(3)
	go c.readLoop()
	go c.handleLoop()
	go c.writeLoop()
}

func (c *Channel) Init() *Channel {
	c.cancelCtx, c.cancel = context.WithCancel(context.TODO())
	c.writer = bufio.NewWriter(c.conn)
	c.reader = bufio.NewReader(c.conn)
	return c
}

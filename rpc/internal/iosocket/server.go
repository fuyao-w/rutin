package iosocket

import (
	"github.com/fuyao-w/rutin/iokit"
)

type Context struct {
	Payload []byte
}

type Server struct {
	Server *iokit.IoServer
}

func NewServer(cb iokit.HandlerFunc) *Server {
	if cb == nil{
		//fmt.Println("newe server cb == nil")
		//os.Exit(1)
	}

	ioServer := iokit.NewIoServer(iokit.NewOnMessage(cb), iokit.NewCodec(&iokit.ProtocolParser{}))
	return &Server{
		Server: ioServer,
	}
}
//
//func (i *IoSocket) ServerOnMessage(body []byte, wc io.WriteCloser) {
//	desc, err := metadata.Parse(i.codec, body)
//	if err != nil {
//		return
//	}
//
//}

package server

import (
	"fmt"
	"github.com/fuyao-w/sd/sd"
)

var DefaultServer Server

func ConfigServer(server Server) {
	server.addr = fmt.Sprintf("%s:%d", "127.0.0.1", server.Port)
	DefaultServer = server
	sd.DefaultRegisterCenter.Register(server.Name, server.addr)
}

package proxy

import (
	"github.com/fuyao-w/sd/proxy/client"
	"github.com/fuyao-w/sd/proxy/server"
)

type Config struct {
	Server  server.Server   `toml:"server"`
	Clients []client.Client `toml:"clients"`
}


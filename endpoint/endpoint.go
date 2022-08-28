package endpoint

import "net"

type ServiceInfo struct {
	DC, Name string
	Tags     []string
	Addr     net.Addr
}

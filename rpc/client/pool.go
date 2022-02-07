package client

import (
	"github.com/jonboulle/clockwork"
	"sync"
)

type dialer interface {
	Dial(host string) (socket, error)
}
type Dial struct {
	options Options
}

func (d *Dial) Dial(host string) (socket, error) {
	return nil, nil
}

type pool struct {
	size    int
	ttl     int64
	cleanup int64
	d       dialer

	// protect sockets
	sync.Mutex
	sockets map[string]map[socket]int64

	clock clockwork.Clock
}

func (p *pool) GetSocket(host string) (socket, error) {
	return p.d.Dial(host)
}

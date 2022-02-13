package client

import (
	"github.com/jonboulle/clockwork"
	"log"
	"net"
	"sync"
	"time"
)

const (
	DefaultPoolSize = 10
	DefaultPoolTTL  = time.Second * 60
)

type dialer interface {
	Dial(host string) (socket, error)
}
type defaultDialer struct {
	options Options
}

func (d *defaultDialer) Dial(host string) (socket, error) {
	conn, err := net.DialTimeout("tcp", host, d.options.Timeout)
	if err != nil {
		log.Printf("defaultDialer|DialTimeout err %s,host :%s ,%v", err, host, d.options.Timeout.Milliseconds())
		return nil, err
	}
	return initRpcSocket(conn, d.options), err
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

func (p *pool) releaseSock(host string, sock socket) {
	p.Lock()
	defer p.Unlock()
	if p.sockets[host] == nil {
		return
	}
	if _, ok := p.sockets[host][sock]; ok {
		return
	}
	delete(p.sockets[host], sock)
	sock.Close()
}
func (p *pool) getSocket(host string) (sock socket, err error) {
	var addSock = func(sock socket) {
		p.sockets[host][sock] = p.clock.Now().Unix()
	}
	defer func() {
		//log.Printf("getsocket size :%d", len(p.sockets[host]))
	}()
	p.Lock()
	defer p.Unlock()
	now := p.clock.Now().Unix()
	if now-p.cleanup > 60 {
		for _, socks := range p.sockets {
			for sock, lastActivity := range socks {
				if now-lastActivity > p.ttl {
					delete(socks, sock)
					go sock.Close()
				}
			}
		}
		p.cleanup = now
	}
	if len(p.sockets[host]) < p.size {
		//p.Unlock()
		sock, err := p.d.Dial(host)
		//p.Lock()
		if err != nil {
			return nil, err
		}
		if p.sockets[host] == nil {
			p.sockets[host] = make(map[socket]int64, 1)
		}

		addSock(sock)
		return sock, err
	}

	for sock, lastActivity := range p.sockets[host] {
		if now-lastActivity > p.ttl {
			continue
		}

		return sock, nil
	}
	//log.Println("add sock",len(p.sockets[host]) , p.size)
	//p.Unlock()
	sock, err = p.d.Dial(host)
	//p.Lock()
	addSock(sock)
	return sock, err
}

func newPool(size int, ttl time.Duration, d dialer) pool {
	clock := clockwork.NewRealClock()
	return pool{
		size: func() int {
			if size > 0 {
				return size
			}
			panic("pool size must gte 0")
		}(),
		ttl: func() int64 {
			if ttl.Seconds() > 0 {
				return int64(ttl.Seconds())
			}
			panic("pool ttl must gte 0")
		}(),
		cleanup: clock.Now().Unix(),
		sockets: map[string]map[socket]int64{},
		d:       d,
		clock:   clock,
	}
}

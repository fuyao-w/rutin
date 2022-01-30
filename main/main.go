package main

import (
	"github.com/fuyao-w/sd/service_discover"
	"time"
)

func main() {
	go service_discover.Consumer()
	time.Sleep(time.Second)
	service_discover.Producer()
}

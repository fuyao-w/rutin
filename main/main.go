package main

import (
	"code.inke.cn/lingxi/server/athena/sd/service_discover"
	"time"
)

func main() {
	go service_discover.Consumer()
	time.Sleep(time.Second)
	service_discover.Producer()
}

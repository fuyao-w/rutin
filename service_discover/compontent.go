package service_discover

import (
	"bufio"
	"fmt"
	"github.com/fuyao-w/sd/net"
	"github.com/fuyao-w/sd/parse"
	netAddr "net"
	"os"
)

func GetIP() string {
	addrs, err := netAddr.InterfaceAddrs()
	if err != nil {
		fmt.Println("get ip err", err)
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*netAddr.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func Consumer() {
	adrr := fmt.Sprintf("%s:10000", GetIP())

	sd, err := NewRedisRegisterProtocol()
	if err != nil {
		fmt.Println("init err ", err)
		return
	}
	sd.SetRegisterName("wfy")
	if err = sd.Register(adrr); err != nil {
		fmt.Println("consumer register ", err)
		return
	}
	net.Server(adrr, net.HandleConnection)
}

func Producer() {
	sd, err := NewRedisRegisterProtocol()
	if err != nil {
		fmt.Println("init err ", err)
		return
	}
	sd.SetRegisterName("wfy")
	arrs := sd.GetAddrSlice()
	if len(arrs) == 0 {
		fmt.Println("GetAddrSlice zero")
	}
	input := bufio.NewScanner(os.Stdin)
	msg := make(chan string)
	go net.Client(arrs[0], msg, &parse.ProtocolParser{})
	fmt.Println("init")
	for input.Scan() {
		line := input.Text()
		// 输入bye时 结束
		if line == "bye" {
			close(msg)
			break
		}
		msg <- line
	}
}

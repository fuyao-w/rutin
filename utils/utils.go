package utils

import (
	"encoding/json"
	netAddr "net"
)

func GetJsonBytes(obj interface{}) []byte {
	bytes, _ := json.Marshal(obj)
	return bytes
}
func GetIP() string {
	addrs, err := netAddr.InterfaceAddrs()
	if err != nil {
		//fmt.Println("get ip err", err)
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
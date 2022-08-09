package utils

import (
	"encoding/json"
	"io"
	netAddr "net"
	"runtime"
)

func PrintLine() (funcName, file string, line int) {
	var (
		ok      bool
		funcPtr uintptr
	)
	funcPtr, file, line, ok = runtime.Caller(1)
	if !ok {
		//fmt.Println("func name: " + )
		//fmt.Printf("file: %s, line: %d\n", file, line)
		return
	}
	if pc := runtime.FuncForPC(funcPtr); pc != nil {
		funcName = pc.Name()
	}

	return
}

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

type MockReader []byte

// MockReader mock Reader interface , r 的长度不要大于 bufio.Reader 缓冲区的长度
func (r MockReader) Read(p []byte) (n int, err error) {
	copy(p, r)
	return len(r), io.EOF
}

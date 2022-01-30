package net

import (
	"bufio"
	"fmt"
	"github.com/fuyao-w/sd/parse"
	"net"
)

/*
	'' ->  '-'
	'//-' -> '///-'

*/
func HandleConnection(conn net.Conn) {
	parser := parse.ProtocolParser{}
	defer conn.Close()
	fmt.Println("handleConnection")

	for {
		bytes, err := parser.Decode(conn)
		fmt.Println(string(bytes), err)
	}

}

func Server(addr string, handle func(conn net.Conn)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("net.Listen err ", err)
		return
	}
	fmt.Println("server init")
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept err", err)
			continue
		}
		fmt.Println("handel")
		go handle(conn)
	}
}

func Client(addr string, msg chan string, parser parse.MsgParser) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial err")
		return
	}

	defer parser.Close()
	for s := range msg {

		w := bufio.NewWriter(conn)
		bytes, _ := parser.Encode([]byte(s))

		_, err := w.Write(bytes)
		if err != nil {
			fmt.Println("put err", err)
		}
		w.Flush()
	}

}

// 分隔符：按字符读取 遇到分隔符判断前面是否有转移符号，如果有就继续读取，如果没有则分割。 全部读取完毕后将 转移符号+ 分割符号 替换成分隔符

package net

import (
	"bufio"
	"fmt"
	"github.com/fuyao-w/sd/parse"
	"log"
	"net"
)

/*
	'' ->  '-'
	'//-' -> '///-'

*/

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

func Client(addr string, msg []byte, parser parse.MsgParser) ([]byte, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("dial err")
		return nil, err
	}

	defer parser.Close()
	w := bufio.NewWriter(conn)
	bytes, _ := parser.Encode(msg)

	_, err = w.Write(bytes)
	if err != nil {
		fmt.Println("put err", err)
	}
	if err = w.Flush(); err != nil {
		log.Printf("client flush err :%s", err.Error())
		return nil, err
	}

	response, err := parser.Decode(conn)
	if err != nil {
		log.Println("Decode err", err)
		return nil, err
	}
	return response, nil
}

// 分隔符：按字符读取 遇到分隔符判断前面是否有转移符号，如果有就继续读取，如果没有则分割。 全部读取完毕后将 转移符号+ 分割符号 替换成分隔符

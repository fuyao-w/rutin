package iokit

import (
	"bufio"
	"io"
)

/*
	自定义分隔符的解析器
	@Delim 自定义的分隔符
	只要消息体里也能出现分隔符的协议都不适合这个解码器
*/
type DelimParser struct {
	Delim byte
}

func (d *DelimParser) Decode(conn *bufio.Reader) (bytes []byte, err error) {
	r := bufio.NewReader(conn)
	for {
		n, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if n == d.Delim {
			break
		} else {
			bytes = append(bytes, n)
		}
	}
	return
}

func (d *DelimParser) Encode(content []byte) (bytes []byte, err error) {
	return append(content, d.Delim), nil
}

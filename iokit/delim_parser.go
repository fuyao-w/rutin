package iokit

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

/*
	自定义分隔符的解析器
	@Delim 自定义的分隔符
	原理：在内容结尾追加分隔符
	如果内容里也有分隔符则用 \ 符号转义，
	解析的时候先根据分隔符取内容，然后再将 内容里的转义符替换回来
 */
type DelimParser struct {
	Delim byte
}


func (d *DelimParser) Encode(content []byte) (bytes []byte, err error) {
	for _, b := range content {
		if b == d.Delim {
			/*
				aaa   -> aaa-
				aaa-  -> aaa /--
				aaa/- -> aaa //--
				aaa/-- -> aaa //-/--
				aaa/\- -> aaa /\/--
			*/
			bytes = append(bytes, '/')
		}

		bytes = append(bytes, b)
	}
	bytes = append(bytes, d.Delim)
	return bytes, nil
}

func (d *DelimParser) Close() {

}

func (d *DelimParser) Decode(conn net.Conn) (bytes []byte, err error) {
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
			if len(bytes) > 0 && bytes[len(bytes)-1] == '/' {
				bytes = append(bytes, n)
			} else {
				delim := string([]byte{d.Delim})
				result := strings.ReplaceAll(string(bytes), fmt.Sprintf("/%s", delim), delim)
				bytes = []byte(result)
				break
			}
		} else {
			bytes = append(bytes, n)
		}
	}

	return
}

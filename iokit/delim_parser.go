package iokit

import (
	"bufio"
	"io"
)

/*
	DelimParser 分隔符解析器，使用 $ 作为分隔符，% 作为转义符
*/
type DelimParser struct{}

var (
	escapes byte = '%'
	delim   byte = '$'
)

func (d *DelimParser) Encode(content []byte) (bytes []byte, err error) {
	for _, b := range content {
		if b == escapes || b == delim {
			bytes = append(bytes, escapes)
		}
		bytes = append(bytes, b)
	}
	bytes = append(bytes, delim)
	return
}

func (d *DelimParser) Decode(r *bufio.Reader) (bytes []byte, err error) {
	var (
		last bool
	)
	for {
		n, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return bytes, err
		}

		if last {
			bytes = bytes[:len(bytes)-1]
			bytes = append(bytes, n)
			last = false
		} else {
			if n == delim {
				break
			}
			bytes = append(bytes, n)
			if n == escapes {
				last = true
			}
		}
	}
	return
}

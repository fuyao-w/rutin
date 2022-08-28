package iokit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type ProtocolParser struct{}

func (p *ProtocolParser) Decode(r *bufio.Reader) ([]byte, error) {
	var (
		fieldMap = map[string]string{}
		key      string
	)
	line, err := r.ReadString('\n')
	if err != nil {
		log.Println("ReadString err ", err)
		return nil, err
	}

	line = strings.TrimRight(line, "\n")
	for i, field := range strings.Split(line, " ") {
		if i&1 == 1 {
			fieldMap[key] = field
		} else {
			key = field
		}
	}
	val, ok := fieldMap["length"]
	if !ok {
		return nil, errors.New("protocol not have field length")
	}
	length, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("parse field length err: %s", err)
	}
	buf := make([]byte, length)
	//这块读取的时候可能包还没发过来呢，所以用 io.ReadFull
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, fmt.Errorf("read content err: %s", err)
	}
	if n != length {
		return nil, errors.New("read content err: length not enough")
	}
	//去掉末尾的\n
	return buf[:len(buf)-1], err
}

/*
协议
header\n
content

header 内容：版本 content 长度（字节数）
version 1.1 length 10093
*/
func (p *ProtocolParser) Encode(content []byte) ([]byte, error) {
	var builder strings.Builder
	content = append(content, '\n')
	builder.WriteString(genHeader(len(content)))
	builder.Write(content)
	return []byte(builder.String()), nil
}

func genHeader(length int) string {
	return fmt.Sprintf("version 0.1 length %d\n", length)
}

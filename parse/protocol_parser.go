package parse

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ProtocolParser struct{}

func (d *ProtocolParser) Clone() MsgParser {
	return &ProtocolParser{}
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

func (p *ProtocolParser) Decode(conn net.Conn) ([]byte, error) {
	var (
		r        = bufio.NewReader(conn)
		fieldMap = map[string]string{}
		key      string
	)
	line, err := r.ReadString('\n')
	if err != nil {
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
	n, err := r.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read content err: %s", err)
	}
	if n != length {
		return nil, errors.New("read content err: length not enough")
	}
	//去掉末尾的\n
	return buf[:len(buf)-1], err
}

func (p *ProtocolParser) Close() {

}

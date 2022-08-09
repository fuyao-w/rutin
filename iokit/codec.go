package iokit

import (
	"bufio"
)

type MsgCodec interface {
	Encode(content []byte) ([]byte, error)
	Decode(conn *bufio.Reader) ([]byte, error)
}

type Packet interface {
	Encode() ([]byte, error)
	Decode(packet []byte) error
}

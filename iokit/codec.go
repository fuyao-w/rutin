package iokit

import "io"

type MsgCodec interface {
	Encode(content []byte) ([]byte, error)
	Decode(conn io.ReadCloser) ([]byte, error)
}

type Packet interface {
	Encode() ([]byte, error)
	Decode(packet []byte) error
}

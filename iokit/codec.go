package iokit

import "io"

type MsgCodec interface {
	Encode(content []byte) ([]byte, error)
	Decode(conn io.ReadCloser) ([]byte, error)
}
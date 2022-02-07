package codec

import (
	"bytes"
	"encoding/gob"
)

type GobCodec struct {
}

func (g *GobCodec) Encode(i interface{}) (buf []byte, err error) {
	err = gob.NewEncoder(bytes.NewBuffer(buf)).Encode(i)
	return buf, err
}

func (g *GobCodec) Decode(buf []byte, i interface{}) error {
	return gob.NewDecoder(bytes.NewReader(buf)).Decode(i)
}

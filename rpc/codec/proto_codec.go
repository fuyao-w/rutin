package codec

import "github.com/gogo/protobuf/proto"

type ProtoCodec struct {
}

func (p *ProtoCodec) Encode(i interface{}) ([]byte, error) {
	return proto.Marshal(i.(proto.Message))
}

func (p *ProtoCodec) Decode(bytes []byte, i interface{}) error {
	return proto.Unmarshal(bytes,i.(proto.Message))
}

package metadata

import "github.com/fuyao-w/rutin/rpc/codec"

type HandlerDesc struct {
	ServiceName string `json:"service_name"`
	MethName    string `json:"meth_name"`
	Param       []byte `json:"param"`
	Response    []byte `json:"response"`
	SeqID       uint64 `json:"seq_id"`
}

func Parse(codec codec.RequestCodec, body []byte) (desc HandlerDesc, err error) {
	err = codec.Decode(body, &desc)
	return
}

package metadata

import (
	"encoding/json"
)

type HandlerDesc struct {
	ServiceName string `json:"service_name"`
	MethName    string `json:"meth_name"`
	Param       []byte `json:"param"`
	Response    []byte `json:"response"`
}

func Unmarshal(body []byte) (desc HandlerDesc, err error) {
	err = json.Unmarshal(body, &desc)
	return
}

func Marshal(desc *HandlerDesc) (body []byte, err error) {
	return json.Marshal(desc)
}

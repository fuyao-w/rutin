package codec

import "encoding/json"

type JsonCodec struct {
}

func (j *JsonCodec) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *JsonCodec) Decode(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}

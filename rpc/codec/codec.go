package codec

//请求协议
type RequestCodec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

package iokit
/*
	codec.go 通过两个层次对网络包进行编、解码
	第一层：基本的拆包、封包，通过 MsgCodec 来实现
	第二层：在一个完整网络请求包基础上对上层应用协议进行封装，进而满足长链接条件下的消息顺序问题

	通过两层协议区分是的整体传输流程比较容易理解，但是代码实现稍复杂。如果同时使用，每个请求至少各编解码两次

	如果使用桥接的方式：
	type MsgCodec interface {
		Encode(content []byte) (Packet, error)
		Decode(pck Packet) ([]byte, error)
	}
	上面实现可以减少调用处代码复杂度，但是依然需要硬编码将两种编码方式组合起来，导致不容易理解。
	这里以容易学习理解为主，所以采用两层编解码的方式

 */
import "io"

type MsgCodec interface {
	Encode(content []byte) ([]byte, error)
	Decode(conn io.ReadCloser) ([]byte, error)
}

type Packet interface {
	Encode() ([]byte, error)
	Decode(packet []byte) error
}
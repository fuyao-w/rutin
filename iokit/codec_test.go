package iokit

import (
	"testing"
)

type read []byte

func TestName(t *testing.T) {

	//var bytes []byte
	d := DelimParser{}

	res, err := d.EncodeC([]byte("%$$%$$$"))
	t.Log(string(res), err)
	////abytes, err := d.Encode([]byte("aaa/"))
	////t.Log(string(abytes), err)
	//abytes := []byte("aaa//\\-/-")
	//for _, n := range abytes {
	//
	//	if n == d.Delim {
	//		_ = "sfdsf\r\nsdfsdf"
	//		if len(bytes) > 0 && bytes[len(bytes)-1] == '/' {
	//			bytes = append(bytes, n)
	//		} else {
	//			delim := string([]byte{d.Delim})
	//			result := strings.ReplaceAll(string(bytes), fmt.Sprintf("/%s", delim), delim)
	//			bytes = []byte(result)
	//			break
	//		}
	//	} else {
	//		bytes = append(bytes, n)
	//	}
	//}
	//t.Log(string(bytes))

}

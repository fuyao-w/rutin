package iokit

import (
	"bufio"
	"github.com/fuyao-w/rutin/utils"
	"testing"
)

var tests = []string{
	"%$$%$$$",
	"dd4$%$%",
	"%%%$123iji^^$(",
	"!232144555",
	"1blsdfiglkajdgfiaifa",
	`{"name":"jack"}"`,
	`\n/b\t\z/nbf/\\/.\`,
	"00|||||`````c..",
}

func TestDelimParser(t *testing.T) {
	testParser(t, &DelimParser{})
}

func TestProtocolParser(t *testing.T) {
	testParser(t, &ProtocolParser{})
}

func testParser(t *testing.T, d MsgCodec) {
	for _, test := range tests {
		res, _ := d.Encode([]byte(test))
		bytes, _ := d.Decode(bufio.NewReader(utils.MockReader(res)))
		if test != string(bytes) {
			t.Log(test, "fail ->", string(bytes))
			t.FailNow()
		}
	}
}

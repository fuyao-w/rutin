package iokit

import "encoding/json"

type SeqPacket struct {
	SeqID   uint64 `json:"seq_id"`
	Payload []byte `json:"payload"`
}

func (s *SeqPacket) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SeqPacket) Decode(packet []byte) error {
	return json.Unmarshal(packet, &s)
}

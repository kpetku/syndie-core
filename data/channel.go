package data

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type Channel struct {
	Name          string
	EncryptKey    string
	Identity      string
	IdentHash     string
	ReadKeys      string
	PublicReplies bool
	Edition       int
	Description   string
	Avatar        []byte
}

func (c *Channel) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := msgpack.NewEncoder(buf)
	err := enc.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Channel) Decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := msgpack.NewDecoder(buf)
	err := dec.Decode(&c)
	if err != nil {
		return err
	}
	return nil
}

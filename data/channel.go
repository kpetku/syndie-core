package data

import (
	"bytes"
	"encoding/gob"
)

type Channel struct {
	Name          string
	EncryptKey    string
	Identity      string
	IdentHash     string
	ReadKeys      []string
	PublicReplies bool
	Edition       int
	Description   string
	Avatar        []byte
}

func (c *Channel) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Channel) Decode(data []byte) (*Channel, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func DecodeChannel(c *Channel, data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&c)
	return err
}

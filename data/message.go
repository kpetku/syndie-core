package data

import (
	"bytes"
	"encoding/gob"
)

type Message struct {
	ID            int
	Name          string
	Subject       string
	Body          string
	TargetChannel string
}

func (m *Message) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) Decode(data []byte) (*Message, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func DecodeMessage(m *Message, data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&m)
	return err
}

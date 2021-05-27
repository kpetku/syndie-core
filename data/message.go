package data

import (
	"bytes"
	"encoding/gob"

	"github.com/kpetku/libsyndie/syndieutil"
)

type Message struct {
	ID            int
	Subject       string
	Body          string
	TargetChannel string
	Avatar        []byte
	Author        string
	PostURI       syndieutil.URI
	Raw           syndieutil.Message
	Header        syndieutil.Header
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

func (m *Message) Decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&m)
	if err != nil {
		return err
	}
	return nil
}

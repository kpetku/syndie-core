package archive

import (
	"encoding/binary"
	"errors"
	"io"
	"log"

	"github.com/hkparker/go-i2p/lib/common/base64"
)

const upperBoundLimit = 10000

type reader struct {
	r   io.Reader
	err error
}

func (r *reader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

type ArchiveHeader struct {
	ArchiveFlags uint16
	AdminChannel uint32
	AltURIs      []string
	NumAltURIs   byte
	NumChannels  uint32
	numMessages  uint32
}

type ArchiveChannelHash struct {
	ChannelHash    [32]byte
	ChannelEdition uint64
	ChannelFlags   byte
}

type ArchiveMessage struct {
	MessageID     uint64
	ScopeChannel  uint32
	TargetChannel uint32
	MsgFlags      byte
}

func Parse(input io.Reader) (*ArchiveHeader, error) {
	var ah ArchiveHeader
	var archiveChannelHashes []ArchiveChannelHash
	var archiveMessages []ArchiveMessage
	var archiveAltURIs []string

	r := reader{r: input}

	// Read ArchiveFlags (unimplemented)
	r.read(&ah.ArchiveFlags)

	// Read the admin channel
	r.read(&ah.AdminChannel)

	// Count the number of alternate URIs
	r.read(&ah.NumAltURIs)
	if int(ah.NumAltURIs) > upperBoundLimit {
		return &ah, errors.New("invalid syndie archive server, too many alternate archive URIs")
	}

	// Populate AltURIs with other known archive servers
	for i := 0; i < int(ah.NumAltURIs); i++ {
		var length uint16
		r.read(&length)

		uri := make([]byte, int(length))
		r.read(&uri)

		archiveAltURIs = append(archiveAltURIs, string(uri))
	}
	ah.AltURIs = archiveAltURIs

	// Count the number of channels
	r.read(&ah.NumChannels)
	if int(ah.NumChannels) > upperBoundLimit {
		return &ah, errors.New("invalid syndie archive server, too many channels")
	}

	// Read the channel hashes
	for i := 0; i < int(ah.NumChannels); i++ {
		var hash ArchiveChannelHash
		r.read(&hash)
		archiveChannelHashes = append(archiveChannelHashes, hash)
	}

	r.read(&ah.numMessages)
	if int(ah.numMessages) > upperBoundLimit {
		return &ah, errors.New("invalid syndie archive server, too many messages")
	}

	// Read messages
	// TODO: lots
	for i := 0; i < int(ah.numMessages); i++ {
		var message ArchiveMessage
		r.read(&message)
		log.Printf("Found messageID: %d, target: %s, scope, %s\n", int(message.MessageID), base64.I2PEncoding.EncodeToString(archiveChannelHashes[int(message.TargetChannel)].ChannelHash[:]), base64.I2PEncoding.EncodeToString(archiveChannelHashes[int(message.ScopeChannel)].ChannelHash[:]))
		archiveMessages = append(archiveMessages, message)
	}

	if r.err != nil {
		return nil, r.err
	}

	return &ah, nil
}

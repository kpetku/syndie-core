package fetcher

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/kpetku/go-syndie/syndieutil"
	"github.com/kpetku/syndied/data"
)

func ImportFile(name string) error {
	dat, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	outer := syndieutil.New()
	_, oerr := outer.Unmarshal(bytes.NewReader(dat))
	if oerr != nil {
		return oerr
	}
	if outer.MessageType == "meta" {
		c := data.Channel{}
		c.Identity = outer.Identity
		c.Edition = outer.Edition
		c.EncryptKey = outer.EncryptKey
		c.Name = outer.Name
		c.Description = outer.Description
		c.ReadKeys = outer.ChannelReadKeys

		encoded, _ := c.Encode()
		foo, _ := syndieutil.ChanHash(c.Identity)

		data.WriteChannel([]byte(foo), encoded)
	}
	if outer.MessageType == "post" {
		out := data.Message{}
		outer := syndieutil.New()
		_, err4 := outer.Unmarshal(bytes.NewReader(dat))
		if err4 != nil {

		}
		lookup, err := data.ReadChannel([]byte(outer.TargetChannel))
		if lookup == nil || err != nil {
			return err
		}
		if lookup.ReadKeys[0] != "" {
			inner := syndieutil.New(syndieutil.BodyKey(lookup.ReadKeys[0]))
			dat2, _ := ioutil.ReadFile(name)
			msg, err3 := inner.Unmarshal(bytes.NewReader(dat2))
			if err3 != nil {
				return err3
			}
			if inner.Subject != "" {
				log.Printf("Subject: %s", inner.Subject)
				log.Printf("Body: %s", msg.Page[0].Data)
			}
			out.Author = inner.Author
			out.TargetChannel = inner.TargetChannel
			out.Avatar = msg.Avatar
			out.Name = inner.Name
			out.Subject = inner.Subject
			out.Raw = msg
			encoded, _ := out.Encode()
			data.WriteMessage([]byte(strconv.Itoa(inner.PostURI.MessageID)), encoded)
		}
	}
	return nil
}

func FetchFromDisk(path string) {
	fetchChannelList, _ := ioutil.ReadDir(path)
	for _, c := range fetchChannelList {
		if c.IsDir() {
			FetchFromDisk(path + c.Name())
			continue
		}
		err := ImportFile(path + "/" + c.Name())
		if err != nil {
			continue
		}
	}
}

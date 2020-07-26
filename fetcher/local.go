package fetcher

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kpetku/libsyndie/syndieutil"
	"github.com/kpetku/syndie-core/data"
)

// LocalFile opens a file from the path location and imports it into the database
func (f *Fetcher) LocalFile(location string) error {
	dat, err := ioutil.ReadFile(location)
	if err != nil {
		return err
	}
	outer := syndieutil.New()
	inner, _ := outer.Unmarshal(bytes.NewReader(dat)) // intentionally ignore the error
	if outer.MessageType == "meta" {
		hidden := syndieutil.New(syndieutil.BodyKey(outer.BodyKey))
		c := data.Channel{}
		c.Identity = outer.Identity
		c.Edition = outer.Edition
		c.EncryptKey = outer.EncryptKey
		c.Name = outer.Name
		c.Description = outer.Description
		c.ReadKeys = outer.ChannelReadKeys + " " + hidden.ChannelReadKeys

		if inner != nil {
			if len(inner.Avatar) > 0 {
				c.Avatar = inner.Avatar
			}
		}
		encoded, errencoding := c.Encode()
		if errencoding != nil {
			log.Printf("errencoding err: %s", errencoding)
		}

		foo, err := syndieutil.ChanHash(c.Identity)
		if err != nil {
			log.Printf("Chanhash err: %s", err)
		}

		cerr := data.WriteChannel([]byte(foo), encoded)
		if cerr != nil {
			log.Printf("error in WriteChannel: %s", cerr)
		}
		log.Printf("wrote metadata for file: %s", location)
		return nil
	}
	if outer.MessageType == "post" {
		out := data.Message{}
		outer := syndieutil.New()
		outer.Unmarshal(bytes.NewReader(dat))
		var lookup *data.Channel
		if outer.TargetChannel != "" {
			lookup, err = data.ReadChannel([]byte(outer.TargetChannel))
		} else {
			lookup, err := data.ReadChannel([]byte(outer.PostURI.Channel))
			if lookup == nil || err != nil {
				log.Printf("error reading channel from bolt: %s", err)
				return err
			}
		}
		if lookup == nil {
			return errors.New("No ReadKeys found for message")
		}
		if len(strings.Fields(lookup.ReadKeys)) >= 0 {
			for num, key := range strings.Fields(lookup.ReadKeys) {
				log.Printf("Checking readkey: %d which is: %s", num, key)
				inner := syndieutil.New(syndieutil.BodyKey(key))
				msg, err3 := inner.Unmarshal(bytes.NewReader(dat))
				if err3 != nil {
					continue
				}
				out.Author = inner.Author
				out.TargetChannel = inner.TargetChannel
				out.Avatar = msg.Avatar
				//				out.Name = inner.Name
				out.Subject = inner.Subject
				out.Raw = *msg
				if len(out.Raw.Attachment) > 0 {
					log.Printf("..writing %d attachments to bolt", len(out.Raw.Attachment)-1)
					if len(out.Raw.Attachment[len(msg.Attachment)-1].Data) == 0 {
						log.Printf("inner/msg size is: %d", len(msg.Attachment[len(msg.Attachment)-1].Data))
						log.Printf("outer size is: %d", len(out.Raw.Attachment[len(msg.Attachment)-1].Data))
					}
				}
				if out.Subject == "" {
					out.Subject = "No subject"
				}
				out.PostURI = outer.PostURI
				log.Printf("!!!!!!!!!!!!!out ID is: %d", outer.PostURI.MessageID)
				out.ID = outer.PostURI.MessageID
				encoded, errx := out.Encode()
				if errx != nil {
					log.Printf("error calling Encode: %s", errx)
				}
				erry := data.WriteMessage([]byte(strconv.Itoa(int(out.ID))), encoded)
				if erry != nil {
					log.Printf("error calling WriteMessage: %s", erry)
				}
			}
			return nil
		}
	}
	return nil
}

// LocalDir recursively walks directories of Syndie messages from the path location and imports them into the database
func (f *Fetcher) LocalDir(location string) error {
	return filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		f.LocalFile(path)
		return err
	})
}

package gateway

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kpetku/libsyndie/syndieutil"
	"github.com/kpetku/syndie-core/data"
)

func recentMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><h1>There are %d messages in the db<table>", len(messages))
	fmt.Fprintf(w, "<thead><td>Subject</td><td>Author</td><td>Channel</td><td>Date</td></thead>")
	for _, m := range messages {
		message := syndieutil.URI{}
		message.RefType = "channel"
		message.Channel = m.TargetChannel
		message.MessageID = m.ID
		message.Page = 1

		date := time.Unix(0, int64(message.MessageID)*int64(time.Millisecond))

		author := syndieutil.URI{}
		author.RefType = "channel"
		author.Channel = m.Author

		target := syndieutil.URI{}
		target.RefType = "channel"
		target.Channel = m.TargetChannel

		fmt.Fprintf(w, "<tr><td><a href=\"/"+message.String()+"\">"+m.Subject+"</a></td><td><a href=\"/"+author.String()+"\">"+ChannelName(m.Author)+" "+syndieutil.ShortIdent(m.Author)+"</td><td><a href=\"/"+target.String()+"\">"+ChannelName(target.Channel)+" "+syndieutil.ShortIdent(target.Channel)+"</a></td><td>"+date.Format("2006-01-02")+"</td><tr>")
	}
	fmt.Fprintf(w, "</table></html>")
}

var messages []data.Message

func ChannelName(s string) string {
	var out string
	c := data.Channel{}
	data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte(s)) {
				c.Decode(v)
				out = c.Name
			}
			return nil
		})
		return nil
	})
	if out == "" {
		return "Anonymous"
	}
	return out
}

func PullMessagesFromDB() {
	err := data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		b.ForEach(func(k, v []byte) error {
			m := data.Message{}
			m.Decode(v)
			messages = append(messages, m)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Printf("error in DB view: %s", err)
	}
}

package gateway

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kpetku/libsyndie/syndieutil"
	"github.com/kpetku/syndie-core/data"
	bolt "go.etcd.io/bbolt"
)

func channelsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if path == "" {
		recentMessagesHandler(w, r)
		return
	}
	if path == "channels" {
		log.Printf("Found path: %s", strings.TrimLeft(r.URL.Path, "/"))
		var channels []data.Channel
		w.Header().Set("Content-Type", "text/html")
		data.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("channels"))
			b.ForEach(func(k, v []byte) error {
				c := data.Channel{}
				c.Decode(v)
				channels = append(channels, c)
				return nil
			})
			return nil
		})
		fmt.Fprintf(w, "<html><table>")
		fmt.Fprintf(w, "<thead><td>name</td><td>description</td></thead>")
		for _, c := range channels {
			uri := syndieutil.URI{}
			uri.RefType = "channel"
			foo, _ := syndieutil.ChanHash(c.Identity)
			uri.Channel = foo
			fmt.Fprintf(w, "<tr><td><a href=\"/"+uri.String()+"\">"+c.Name+" "+syndieutil.ShortIdent(foo)+"</a></td><td>"+c.Description+"</td><td>")
		}
		fmt.Fprintf(w, "</table></html>")
		return
	}
	u := syndieutil.URI{}
	err := u.Marshall(path)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html>Error: %s</html>", err.Error())
		return
	}
	if u.MessageID != 0 {
		test, err := getMessage(u.MessageID)
		if err != nil {
			log.Printf("Cannot fetch messageId %d!", u.MessageID)
			return
		}
		date := time.Unix(0, int64(u.MessageID)*int64(time.Millisecond))

		author := syndieutil.URI{}
		author.RefType = "channel"
		author.Channel = test.Author

		target := syndieutil.URI{}
		target.RefType = "channel"
		target.Channel = test.TargetChannel

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><h1>%s</h1>Posted by <a href=\"/%s\">%s</a> in <a href=\"/%s\">%s</a> on %s<p>", test.Subject, author.String(), ChannelName(author.Channel)+" "+syndieutil.ShortIdent(test.Author), target.String(), ChannelName(target.Channel)+" "+syndieutil.ShortIdent(target.Channel), date.Format("2006-01-02"))
		if len(test.Raw.Page) > 0 {
			for num, p := range test.Raw.Page {
				if num >= 0 {
					fmt.Fprintf(w, "<hr><pre>%s</pre>\n", subURI(p.Data))
				}
			}
		}
		if len(test.Raw.Attachment) >= 0 {
			for num, a := range test.Raw.Attachment {
				if len(a.Data) > 0 {
					fmt.Fprintf(w, "attachment number: %d content-type: %s name: %s desc: %s length: %d<p>", num, a.ContentType, a.Name, a.Description, len(a.Data))
					switch filepath.Ext(strings.ToLower(a.Name)) {
					case ".png":
						a.ContentType = "image/png"
					case ".gif":
						a.ContentType = "image/gif"
					case ".jpg", ".jpeg":
						a.ContentType = "image/jpeg"
					}
					fmt.Fprintf(w, `<img src="data:%s;base64,%s" alt="%s">`, a.ContentType, base64.StdEncoding.EncodeToString(a.Data), a.Description)
				}
			}
		}
		fmt.Fprintf(w, "</html>")
		return
	}
	if u.Channel != "" {
		var foundMessages bool
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><table>")
		fmt.Fprintf(w, "<thead><td>subject</td><td>author</td></thead>")
		for _, m := range messages {
			if m.TargetChannel == u.Channel {
				foundMessages = true
				message := syndieutil.URI{}
				message.RefType = "channel"
				message.Channel = m.TargetChannel
				message.MessageID = m.ID
				message.Page = 1

				author := syndieutil.URI{}
				author.RefType = "channel"
				author.Channel = m.Author

				fmt.Fprintf(w, "<tr><td><a href=\"/"+message.String()+"\">"+m.Subject+"</a></td><td><a href=\"/"+author.String()+"\">"+ChannelName(m.Author)+" "+syndieutil.ShortIdent(m.Author)+"</td><tr>")
			}
		}
		fmt.Fprintf(w, "</table></html>")
		if !foundMessages {
			fmt.Fprintf(w, "0 messages found for channel")
		}
		fmt.Fprintf(w, "</html>")
		return
	}
}

func getMessage(i int) (data.Message, error) {
	log.Printf("Checking for messageid: %d", i)
	var out data.Message
	m := data.Message{}
	data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte(strconv.Itoa(i))) {
				m.Decode(v)
				out = m
			}
			return nil
		})
		return nil
	})
	return out, nil
}

func subURI(s string) string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "urn:syndie") {
			log.Printf("subURI called from: %s", scanner.Text())
		}
	}
	return strings.Replace(s, "<img src=\"", "<img src=\"./", -1)
}

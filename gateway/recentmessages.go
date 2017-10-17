package gateway

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/kpetku/syndied/data"
)

func recentMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var messages []data.Message
	w.Header().Set("Content-Type", "text/html")
	data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		b.ForEach(func(k, v []byte) error {
			m := data.Message{}
			m.Decode(v)
			messages = append(messages, m)
			return nil
		})
		return nil
	})
	fmt.Fprintf(w, "<html><table>")
	fmt.Fprintf(w, "<thead><td>subject</td><td>author</td><td>channel</td></thead>")
	for _, m := range messages {
		fmt.Fprintf(w, "<tr><td>"+m.Subject+"</td><td>"+m.Author+"</td><td>"+m.TargetChannel+"</td><tr>")
	}
	fmt.Fprintf(w, "</table></html>")
}

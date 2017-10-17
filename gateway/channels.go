package gateway

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/kpetku/syndied/data"
)

func channelsHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "<tr><td>"+c.Name+"</td><td>"+c.Description+"</td><td>")
	}
	fmt.Fprintf(w, "</table></html>")
}

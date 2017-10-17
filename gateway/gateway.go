package gateway

import (
	"log"
	"net/http"
)

func New() {
	http.HandleFunc("/", channelsHandler)
	http.HandleFunc("/channels", channelsHandler)
	http.HandleFunc("/recentmessages", recentMessagesHandler)
	http.HandleFunc("/rss", rssHandler)
	serr := http.ListenAndServe(":9090", nil)
	if serr != nil {
		log.Fatal("ListenAndServe: ", serr)
	}
}

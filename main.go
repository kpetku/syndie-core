package main

import (
	"flag"
	"log"
	"os/user"

	"github.com/kpetku/syndied/fetcher"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fetchURL := flag.String("fetch", "http://localhost:8080/", "Fetch all known messages from a Syndie archive server")
	fetchPath := flag.String("folder", usr.HomeDir+"/.syndie/incoming", "Specifies which folder to fetch messages into")
	fetchTimeout := flag.Int("timeout", 10, "HTTP timeout value in seconds")
	fetchDelay := flag.Int("delayms", 100, "Impose a random delay of up to n miliseconds when fetching")

	flag.Parse()

	f := fetcher.New(*fetchURL, *fetchPath, *fetchTimeout, *fetchDelay)

	ferr := f.Fetch()
	if ferr != nil {
		log.Printf("Error indexing: %s", err)
	}
}

package main

import (
	"flag"
	"log"
)

const Version string = "v0.1"

func main() {
	fetchURL := flag.String("fetch", "http://syndie.welterde.de", "Fetch all known messages from a Syndie archive server, example: http://syndie.welterde.de")

	flag.Parse()

	f := &Fetcher{}
	_, err := f.Fetch(*fetchURL)
	if err != nil {
		log.Printf("Error fetching: %s", err)
	}
}

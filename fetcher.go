package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kpetku/syndied/archive"
)

type Fetcher struct{}

func (f *Fetcher) Fetch(url string) (*Fetcher, error) {
	log.Printf("syndied version %s fetching: %s", Version, url)

	var c = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", strings.TrimRight(url, "/")+"/shared-index.dat", nil)
	req.Header.Add("User-Agent", "syndied "+Version+" author: keith@keithp.net")
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("Error from resp: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	a, err := archive.Parse(resp.Body)

	if err != nil {
		log.Fatalf("Error from reader: %s", err)
	}
	log.Printf("numAltURIs: %d", a.NumAltURIs)
	log.Printf("NumChannels: %d", int(a.NumChannels))

	return f, nil
}

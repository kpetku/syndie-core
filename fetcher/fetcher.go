package fetcher

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kpetku/go-syndie/archive/client"
)

const upperBoundLimit = 10000

// Fetcher hold things
type Fetcher struct {
	SourceURL  string
	OutputPath string
	Timeout    int
	Batch      int
	Delay      int
	client     *client.Client
}

// New creates a new instance of Fetcher.
func New(remote, path string, timeout, delay int) *Fetcher {
	return &Fetcher{
		SourceURL:  remote,
		OutputPath: path,
		Timeout:    timeout,
		Delay:      delay,
		client:     &client.Client{},
	}
}

// GetIndex reaches out to an endpoint over http and builds a list of urls.
func (f *Fetcher) GetIndex() error {
	req, err := http.NewRequest("GET", strings.TrimRight(f.SourceURL, "/")+"/shared-index.dat", nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "syndied")
	var c = &http.Client{
		Timeout: time.Second * time.Duration(f.Timeout),
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f.client = client.New()
	f.client.Parse(resp.Body)

	log.Printf("numAltURIs: %d", f.client.NumAltURIs)
	log.Printf("NumChannels: %d", int(f.client.NumChannels))
	log.Printf("Number of messages: %d", len(f.client.Urls))

	return nil
}

// Fetch actually fetches all URLs from a remote endpoint into the specified path
func (f *Fetcher) Fetch() error {
	f.GetIndex()
	if f.client.Urls == nil {
		return errors.New("no URLs to fetch")
	}
	if len(f.client.Urls) >= upperBoundLimit {
		return errors.New("too many URLs to fetch")
	}
	for x, url := range f.client.Urls {
		url = strings.TrimRight(f.SourceURL, "/") + "/" + url

		log.Printf("Fetching' %s", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Add("User-Agent", "syndied")
		var c = &http.Client{
			Timeout: time.Second * time.Duration(f.Timeout),
		}
		resp, err := c.Do(req)
		if err != nil {
			resp.Body.Close()
			return err
		}
		defer resp.Body.Close()
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			err := ioutil.WriteFile(f.OutputPath+"/"+strconv.Itoa(1000000000000+rand.Intn(9999999999999-1000000000000))+".syndie", buf, 0644)
			if err != nil {
				return err
			}
			log.Printf("Fetched %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.client.Urls))
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(f.Delay)))
	}
	return nil
}

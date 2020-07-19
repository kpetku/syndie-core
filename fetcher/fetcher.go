package fetcher

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kpetku/libsyndie/archive/client"
	"github.com/kpetku/libsyndie/syndieutil"
)

const upperBoundLimit = 10000

// Fetcher hold things
type Fetcher struct {
	remote    string // remote can be a URL or file
	localPath string // localPath is where to store the results on the local filesystem
	timeout   int    // timeout in second
	delay     int    // delay between individual fetches in miliseconds
	client    *client.Client
}

// New creates a new instance of Fetcher.
func New(remote, path string, timeout, delay int) *Fetcher {
	return &Fetcher{
		remote:    remote,
		localPath: path,
		timeout:   timeout,
		delay:     delay,
		client:    &client.Client{},
	}
}

// GetIndex reaches out to an endpoint over http and builds a list of urls.
func (f *Fetcher) GetIndex() error {
	_, err := url.ParseRequestURI(f.remote)
	if err == nil {
		req, err := http.NewRequest("GET", strings.TrimRight(f.remote, "/")+"/shared-index.dat", nil)
		if err != nil {
			return err
		}
		req.Header.Add("User-Agent", "syndied")
		var c = &http.Client{
			Timeout: time.Second * time.Duration(f.timeout),
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
	fi, err := os.Stat(f.remote)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		fetchChannelList, _ := ioutil.ReadDir(f.remote)
		for _, c := range fetchChannelList {
			if c.IsDir() {
				FetchFromDisk(f.remote + "/" + c.Name())
			} else {
				ImportFile(f.remote + "/" + c.Name())
			}
		}
	} else {
		ImportFile(f.remote)
	}
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
		url = strings.TrimRight(f.remote, "/") + "/" + url

		log.Printf("Fetching' %s", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Add("User-Agent", "syndied")
		var c = &http.Client{
			Timeout: time.Second * time.Duration(f.timeout),
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
			// Validate the message and take the PostURI messageID from it
			outer := syndieutil.New()
			_, err := outer.Unmarshal(bytes.NewReader(buf))
			if err != nil {
				log.Printf("Error unmarshalling outer: %s", err)
			}
			if outer.MessageType == "meta" {
				chanHash, err := syndieutil.ChanHash(outer.Identity)
				if err != nil {
					log.Printf("Error parsing chanhash: %s", err)
				}
				log.Printf("Fetched META %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.client.Urls))
				if _, err := os.Stat(f.localPath + "/" + chanHash + "/"); os.IsNotExist(err) {
					os.Mkdir(f.localPath+"/"+chanHash+"/", 0744)
				}
				dest := f.localPath + "/" + chanHash + "/" + "meta.syndie"
				werr := ioutil.WriteFile(dest, buf, 0644)
				if werr != nil {
					log.Printf("Unable to write post to disk: %s", werr.Error())
				}
				ierr := ImportFile(dest)
				if ierr != nil {
					log.Printf("Unable to import meta: %s", ierr.Error())
				}
				log.Printf("Fetched %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.client.Urls))
			}
			if outer.MessageType == "post" {
				dest := f.localPath + "/" + outer.TargetChannel + "/" + strconv.Itoa(outer.PostURI.MessageID) + ".syndie"
				werr := ioutil.WriteFile(dest, buf, 0644)
				if werr != nil {
					log.Printf("Unable to write post to disk: %s", werr.Error())
				}
				ierr := ImportFile(dest)
				if ierr != nil {
					log.Printf("Unable to import post: %s", ierr.Error())
				}
				log.Printf("Fetched %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.client.Urls))
			}
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(f.delay)))
	}
	return nil
}

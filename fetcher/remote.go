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

// RemoteFetch performs a remote HTTP fetch from "remoteLocation"
func (f *Fetcher) RemoteFetch() error {
	err := f.buildIndex()
	if err != nil {
		f.LocalDir(f.remoteLocation)
	}
	if f.Client.Urls == nil {
		return errors.New("no URLs to fetch")
	}
	if len(f.Client.Urls) >= upperBoundLimit {
		return errors.New("too many URLs to fetch")
	}
	for x, url := range f.Client.Urls {
		url = strings.TrimRight(f.remoteLocation, "/") + "/" + url

		log.Printf("Fetching' %s", url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		tr, err := f.SelectTransport(req.URL.Hostname())
		if err != nil {
			return err
		}
		req.Header.Add("User-Agent", "syndie-core")
		var c = &http.Client{
			Timeout:   time.Second * time.Duration(f.timeout),
			Transport: tr,
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
		log.Printf("Status: %d", resp.StatusCode)
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
				log.Printf("Fetched META %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.Client.Urls))
				if _, err := os.Stat(f.localLocation + "/" + chanHash + "/"); os.IsNotExist(err) {
					os.Mkdir(f.localLocation+"/"+chanHash+"/", 0744)
				}
				dest := f.localLocation + "/" + chanHash + "/" + "meta.syndie"
				werr := ioutil.WriteFile(dest, buf, 0644)
				if werr != nil {
					log.Printf("Unable to write post to disk: %s", werr.Error())
				}
				err = f.LocalFile(dest)
				if err != nil {
					log.Printf("Unable to import meta: %s", err.Error())
				}
				log.Printf("Fetched %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.Client.Urls))
			}
			if outer.MessageType == "post" {
				dest := f.localLocation + "/" + outer.TargetChannel + "/" + strconv.Itoa(outer.PostURI.MessageID) + ".syndie"
				werr := ioutil.WriteFile(dest, buf, 0644)
				if werr != nil {
					log.Printf("Unable to write post to disk: %s", werr.Error())
				}
				ierr := f.LocalFile(dest)
				if ierr != nil {
					log.Printf("Unable to import post: %s", ierr.Error())
				}
				log.Printf("Fetched %s with %d bytes, number: %d/%d", url, len(buf), x, len(f.Client.Urls))
			}
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(f.delay)))
	}
	return nil
}

// buildIndex reaches out to an endpoint over http and builds a list of urls
func (f *Fetcher) buildIndex() error {
	// Try to build the index of a remote archive over HTTP
	_, err := url.ParseRequestURI(f.remoteLocation)
	if err == nil {
		req, err := http.NewRequest("GET", strings.TrimRight(f.remoteLocation, "/")+"/shared-index.dat", nil)
		if err != nil {
			return err
		}
		req.Header.Add("User-Agent", "syndie-core")
		tr, err := f.SelectTransport(req.URL.Hostname())
		if err != nil {
			return err
		}
		var c = &http.Client{
			Timeout:   time.Second * time.Duration(f.timeout),
			Transport: tr,
		}

		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		f.Client = client.New()
		return f.Client.Parse(resp.Body)
	}
	return err
}

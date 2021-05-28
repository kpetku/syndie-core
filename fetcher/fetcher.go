package fetcher

import (
	"log"
	"net/http"
	"strings"

	"github.com/cretz/bine/tor"
	"github.com/eyedeekay/sam3/helper"
	"github.com/eyedeekay/sam3"
	"github.com/kpetku/libsyndie/archive/client"
)

const upperBoundLimit = 10000

// Fetcher contains verious options for a Syndie fetch operation
type Fetcher struct {
	torClient      *tor.Tor
	samClient      *sam3.StreamSession
	remoteLocation string // remoteLocation is a URL pointing to an archive server
	localLocation  string // localLocation is where to store the results on the local filesystem
	timeout        int    // timeout in seconds
	delay          int    // random delay of up to "delay" miliseconds between individual http requests
	SAMAPIaddr     string
	TORSocksaddr   string
	Client         *client.Client
}

// SelectTransport decides, based on the hostname, whether to use an I2P
// or Tor transports
func (f *Fetcher) SelectTransport(host string) (*http.Transport, error) {
	var err error
	if strings.HasSuffix(host, "i2p") {
		if f.SAMAPIaddr != "" {
			//			var sam *goSam.Client
			if f.samClient == nil {
				f.samClient, err = sam.I2PStreamSession("syndie",f.SAMAPIaddr,"syndie") //f.SAMAPIaddr)
				if err != nil {
					return nil, err
				}
			}
			// create a transport that uses SAM to dial TCP Connections
			tr := http.Transport{
				Dial: f.samClient.Dial,
			}
			log.Printf("SAM Client Created")
			return &tr, err
		}
	}
	if f.TORSocksaddr != "" {
		if f.samClient == nil {
			f.torClient, err = tor.Start(nil, nil)
			if err != nil {
				return nil, err
			}
		}
		//defer t.Close()
		//		var dialer *tor.Dialer
		dialer, err := f.torClient.Dialer(nil, nil)
		if err != nil {
			return nil, err
		}
		tr := http.Transport{
			Dial: dialer.Dial,
		}
		log.Printf("Tor Client Created")
		return &tr, err
	}
	log.Printf("Non-Anonymous Client Created")
	return nil, nil
}

// New creates a new instance of Fetcher
func New(remote, path string, timeout, delay int) *Fetcher {
	return &Fetcher{
		remoteLocation: remote,
		localLocation:  path,
		timeout:        timeout,
		delay:          delay,
		Client:         &client.Client{},
		SAMAPIaddr:     "127.0.0.1:7656",
		TORSocksaddr:   "127.0.0.1:9050",
	}
}

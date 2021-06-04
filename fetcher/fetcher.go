package fetcher

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cretz/bine/tor"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/helper"
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
	anonOnly       bool   // only establish anonymous remote connections. Use Tor for all Non-I2P URLs.
	SAMAPIaddr     string // SAM API address to use for setting up connections to archives on I2P
	TORSocksaddr   string // TOR SOCKS Proxy to use for fetching .onion URL's and clearnet URL's anonymously
	Client         *client.Client
}

// SelectTransport decides, based on the hostname, whether to use an I2P
// or Tor transports
func (f *Fetcher) SelectTransport(host string) (*http.Transport, error) {
	var err error
	if strings.HasSuffix(host, "i2p") {
		if f.SAMAPIaddr != "" {
			if f.samClient == nil {
				f.samClient, err = sam.I2PStreamSession("syndie", f.SAMAPIaddr, filepath.Join(f.localLocation, "syndie"))
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
		if strings.HasSuffix(host, "onion") || f.anonOnly {
			if f.torClient == nil {
				f.torClient, err = tor.Start(nil, nil)
				if err != nil {
					return nil, err
				}
			}
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
	}
	if f.anonOnly {
		return nil, errors.New("tor client not available and anonymous-only specified. Failing before we de-anon")
	}
	log.Printf("Non-Anonymous Client Created")
	return &http.Transport{}, nil
}

// NewOpts creates a new instance of Fetcher using a collection of functional
// arguments
func NewOpts(opts ...func(*Fetcher) error) (*Fetcher, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	var f Fetcher
	f.remoteLocation = "http://127.0.0.1:8080"
	f.localLocation = filepath.Join(usr.HomeDir, "/.syndie/incoming")
	f.timeout = 10
	f.delay = 100
	f.anonOnly = false
	f.Client = &client.Client{}
	f.SAMAPIaddr = "127.0.0.1:7656"
	f.TORSocksaddr = "127.0.0.1:9050"
	for _, o := range opts {
		if err := o(&f); err != nil {
			return nil, err
		}
	}
	os.MkdirAll(f.localLocation, 0755)
	return &f, nil
}

// New creates a new instance of Fetcher
func New(remote, path string, timeout, delay int) *Fetcher {
	f, err := NewOpts(
		SetLocalLocation(path),
		SetRemoteLocation(remote),
		SetTimeout(timeout),
		SetDelay(delay),
		SetSAMAPIAddr(""),
		SetTORSocksaddr(""),
	)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// New creates a new instance of Fetcher which will fail if it cannot
// make an anonymous connection.
func NewAnonOnly(remote, path string, timeout, delay int) *Fetcher {
	f, err := NewOpts(
		SetLocalLocation(path),
		SetRemoteLocation(remote),
		SetTimeout(timeout),
		SetDelay(delay),
		SetAnonOnly(true),
		SetSAMAPIAddr("127.0.0.1:7656"),
		SetTORSocksaddr("127.0.0.1:9050"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

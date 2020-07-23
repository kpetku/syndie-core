package fetcher

import (
	"github.com/kpetku/libsyndie/archive/client"
)

const upperBoundLimit = 10000

// Fetcher contains verious options for a Syndie fetch operation
type Fetcher struct {
	remoteLocation string // remoteLocation is a URL pointing to an archive server
	localLocation  string // localLocation is where to store the results on the local filesystem
	timeout        int    // timeout in seconds
	delay          int    // random delay of up to "delay" miliseconds between individual http requests
	Client         *client.Client
}

// New creates a new instance of Fetcher
func New(remote, path string, timeout, delay int) *Fetcher {
	return &Fetcher{
		remoteLocation: remote,
		localLocation:  path,
		timeout:        timeout,
		delay:          delay,
		Client:         &client.Client{},
	}
}

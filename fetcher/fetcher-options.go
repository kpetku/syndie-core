package fetcher

//SetLocalLocation configures the location where the fetcher looks for local
//archives and downloads local archives
func SetLocalLocation(s string) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.localLocation = s
		return nil
	}
}

//SetRemoteLocation configures the remote location to fetch new archives from
func SetRemoteLocation(s string) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.remoteLocation = s
		return nil
	}
}

//SetTimeout configures the timeout
func SetTimeout(s int) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.timeout = s
		return nil
	}
}

//SetAnonOnly tells the Fetcher it should fail if it cannot make an anonymous
//connection
func SetAnonOnly(s bool) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.anonOnly = s
		return nil
	}
}

//SetDelay tells the fetcher the max value for the randomm delay
func SetDelay(s int) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.delay = s
		return nil
	}
}

//SetSAMAPIAddr sets the SAM API address for transparently handling .b32.i2p and .i2p
// archives
func SetSAMAPIAddr(s string) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.SAMAPIaddr = s
		return nil
	}
}

//SetTORSocksaddr sets the TOR SOCKS5 address for transparently handling .onion archives
//and obfuscating the location of clients accessing the clearnet archive
func SetTORSocksaddr(s string) func(*Fetcher) error {
	return func(c *Fetcher) error {
		c.TORSocksaddr = s
		return nil
	}
}

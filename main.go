package main

import (
	"flag"
	"log"
	"os/user"
	"time"

	"github.com/kpetku/syndie-core/data"
	"github.com/kpetku/syndie-core/fetcher"
	"github.com/kpetku/syndie-core/gateway"
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
		log.Printf("Error indexing: %s", ferr)
	}

	derr := data.OpenDB(usr.HomeDir + "/.syndie/db/bolt.db")
	if derr != nil {
		log.Fatal(err)
	}
	defer func() {
		err := data.DB.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	err = data.InitDB()
	if err != nil {
		log.Printf("err: %s", err)
	}

	go gateway.New()
	time.Sleep(time.Second * 60)
	log.Printf("Importing messages from incoming folder to http://localhost:9090/recentmessages")
	fetcher.FetchFromDisk(usr.HomeDir + "/.syndie/incoming/")
	log.Printf("Sleeping for 5 minutes then exiting")
	time.Sleep(time.Minute * 5)
}

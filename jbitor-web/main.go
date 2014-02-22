package main

import (
	"github.com/jbitor/dht"
	"github.com/jbitor/webclient"
	"log"
	weakrand "math/rand"
	"os"
	"time"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", 0)
}

func main() {
	if len(os.Args) == 0 {
		logger.Fatalf("Usage: %v\n", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	dc, err := dht.OpenClient(".dht-peer", false)
	if err != nil {
		logger.Fatalf("Unable to open DHT client: %v\n", err)
		return
	}

	wc, err := webclient.NewForDhtClient(dc)
	if err != nil {
		logger.Fatalf("Unable to create web client: %v\n", err)
		return
	}

	err = wc.ListenAndServe()
	if err != nil {
		logger.Fatalf("Unable to serve web client: %v\n", err)
		return
	}

	defer dc.Close()

	for {
		time.Sleep(60 * time.Second)
	}
}

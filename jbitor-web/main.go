package main

import (
	weakrand "math/rand"
	"os"
	"time"

	"github.com/jbitor/bittorrent"
	"github.com/jbitor/bittorrent/dht"
	"github.com/jbitor/cli/loggerconfig"
	"github.com/jbitor/webclient"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func main() {
	loggerconfig.Use()

	if len(os.Args) == 0 {
		logger.Fatalf("Usage: %v", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	dc, err := dht.OpenClient(".dht-peer", false)
	if err != nil {
		logger.Fatalf("Unable to open DHT client: %v", err)
		return
	}

	bc := bittorrent.OpenClient()

	wc, err := webclient.New(dc, bc)
	if err != nil {
		logger.Fatalf("Unable to create web interface: %v", err)
		return
	}

	err = wc.ListenAndServe()
	if err != nil {
		logger.Fatalf("Unable to serve web interface: %v", err)
		return
	}

	defer dc.Close()

	for {
		time.Sleep(60 * time.Second)
	}
}

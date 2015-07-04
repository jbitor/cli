package main

import (
	weakrand "math/rand"
	"os"
	"time"

	"github.com/jbitor/dht"
	"github.com/jbitor/webclient"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func init() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{level:4.4s} %{id:03x}%{color:reset} %{message}\n         %{longfunc}() in %{module}/%{shortfile}\n\n",
		)))
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

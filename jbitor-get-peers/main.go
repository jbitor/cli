package main

import (
	"encoding/json"

	weakrand "math/rand"
	"os"
	"time"

	"github.com/op/go-logging"

	"github.com/jbitor/bittorrent"
	"github.com/jbitor/dht"
)

var logger = logging.MustGetLogger("main")

func init() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{level:4.4s} %{id:03x}%{color:reset} %{message}\n         %{longfunc}() in %{module}/%{shortfile}\n\n",
		)))
}

func main() {
	if len(os.Args) != 2 {
		logger.Fatalf("Usage: %v INFOHASH\n", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	infoHash, err := bittorrent.BTIDFromHex(os.Args[1])

	if err != nil {
		logger.Fatalf("Specified string was not a valid hex infohash [%v].\n", err)
		return
	}

	dhtClient, err := dht.OpenClient(".dht-peer", true)
	if err != nil {
		logger.Fatalf("Unable to open .dht-peer: %v\n", err)
		return
	}

	defer dhtClient.Close()

	search := dhtClient.GetPeers(infoHash)
	peers, err := search.AllPeers()
	if err != nil {
		logger.Fatalf("Unable to find peers: %v\n", err)
	}

	logger.Info("Found peers for %v:\n", infoHash)
	peerData, err := json.Marshal(peers)
	if err != nil {
		logger.Fatalf("?!?: %v\n", err)
	}

	os.Stdout.Write(peerData)
}

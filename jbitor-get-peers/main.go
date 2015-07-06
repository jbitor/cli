package main

import (
	"encoding/json"

	weakrand "math/rand"
	"os"
	"time"

	"github.com/op/go-logging"

	"github.com/jbitor/bittorrent"
	"github.com/jbitor/bittorrent/dht"
	"github.com/jbitor/cli/loggerconfig"
)

var logger = logging.MustGetLogger("main")

func main() {
	loggerconfig.Use()

	if len(os.Args) != 2 {
		logger.Fatalf("Usage: %v INFOHASH", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	infoHash, err := bittorrent.BTIDFromHex(os.Args[1])

	if err != nil {
		logger.Fatalf("Specified string was not a valid hex infohash [%v].", err)
		return
	}

	dhtClient, err := dht.OpenClient(".dht-peer", true)
	if err != nil {
		logger.Fatalf("Unable to open .dht-peer: %v", err)
		return
	}

	defer dhtClient.Close()

	search := dhtClient.GetPeers(infoHash)
	peers, err := search.AllPeers()
	if err != nil {
		logger.Fatalf("Unable to find peers: %v", err)
	}

	logger.Info("Found peers for %v:", infoHash)
	peerData, err := json.Marshal(peers)
	if err != nil {
		logger.Fatalf("?!?: %v", err)
	}

	os.Stdout.Write(peerData)
}

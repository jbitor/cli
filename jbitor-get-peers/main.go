package main

import (
	"fmt"
	"log"
	weakrand "math/rand"
	"os"
	"time"

	"github.com/jbitor/bittorrent"
	"github.com/jbitor/dht"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", 0)
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

	logger.Printf("Found peers for %v:\n", infoHash)
	for _, peer := range peers {
		fmt.Println(peer)
	}
}

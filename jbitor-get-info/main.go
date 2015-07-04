// Finds peers for a torrent and downloads its metadata, saving it as .torrent
// file, with nearby DHT nodes included per BEP-8.
package main

import (
	"encoding/json"
	"log"
	weakrand "math/rand"
	"net"
	"os"
	"time"

	"github.com/jbitor/bittorrent"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", 0)
}

func main() {
	if len(os.Args) != 2 {
		logger.Fatalf("Usage: %v INFOHASH < INFOHASH.peers\n", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	infoHash, err := bittorrent.BTIDFromHex(os.Args[1])

	if err != nil {
		logger.Fatalf("Specified string was not a valid hex infohash [%v].\n", err)
		return
	}

	peers := make([]net.TCPAddr, 0)
	dec := json.NewDecoder(os.Stdin)
	dec.Decode(&peers)
	logger.Printf("Loaded peers: %v\n", peers)

	client := bittorrent.OpenClient()
	swarm := client.Swarm(infoHash)

	for _, peer := range peers {
		swarm.AddPeer(peer)
	}

	info := swarm.Info()

	logger.Printf("got info: %v", info)
}

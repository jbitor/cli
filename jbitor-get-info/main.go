// Finds peers for a torrent and downloads its metadata, saving it as .torrent
// file, with nearby DHT nodes included per BEP-8.
package main

import (
	"encoding/json"
	weakrand "math/rand"
	"net"
	"os"
	"time"

	"github.com/jbitor/bencoding"
	"github.com/jbitor/bittorrent"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func init() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{level:4.4s}%{color:reset} %{message}\n%{color}%{id:4.4x}%{color:reset} %{module} / %{shortfile} / %{longfunc}()\n\n",
		)))
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
	logger.Info("Loaded peers: %v\n", peers)

	client := bittorrent.OpenClient()
	swarm := client.Swarm(infoHash)

	for _, peer := range peers {
		swarm.AddPeer(peer)
	}

	logger.Info("getting info")
	info := swarm.Info()
	logger.Info("got info: %v", info)

	torrentFileData, err := bencoding.Encode(bencoding.Dict{
		"info":          info,
		"announce-list": bencoding.List{},
		"nodes":         bencoding.List{},
	})
	if err != nil {
		logger.Fatalf("error encoding torrent file: %v:", err)
	}

	os.Stdout.Write(torrentFileData)
}

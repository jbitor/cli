package cli

import (
	"fmt"
	"github.com/jbitor/bittorrent"
	"github.com/jbitor/dht"
	"github.com/jbitor/webclient"
	"os"
	"time"
)

func cmdDht(args []string) {
	if len(args) == 0 {
		logger.Fatalf("Usage: %v dht SUBCOMMAND\n", os.Args[0])
		return
	}

	subcommand := args[0]
	subcommandArgs := args[1:]

	switch subcommand {
	case "connect":
		cmdDhtConnect(subcommandArgs)
	case "get-peers":
		cmdDhtGetPeers(subcommandArgs)
	default:
		logger.Fatalf("Unknown dht subcommand: %v\n", subcommand)
		return
	}
}

func cmdDhtConnect(args []string) {
	if len(args) != 0 {
		logger.Fatalf("Usage: %v dht connect.benc\n", os.Args[0])
		return
	}

	client, err := dht.OpenClient(".dht-peer", false)
	if err != nil {
		logger.Fatalf("Unable to open client: %v\n", err)
		return
	}

	err = webclient.ServeForDhtClient(client)
	if err != nil {
		logger.Fatalf("Unable to serve web client: %v\n", err)
		return
	}
	defer client.Close()

	for {
		time.Sleep(60 * time.Second)
	}
}

func cmdDhtGetPeers(args []string) {
	if len(args) != 1 {
		logger.Fatalf("Usage: %v torrent get-peers INFOHASH\n", os.Args[0])
		return
	}

	infoHash, err := bittorrent.BTIDFromHex(args[0])

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

	peers, err := dhtClient.GetPeers(infoHash)

	if err != nil {
		logger.Fatalf("Unable to find peers: %v\n", err)
	}

	logger.Printf("Found peers for %v:\n", infoHash)
	for _, peer := range peers {
		fmt.Println(peer)
	}
}

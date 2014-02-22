package main

import (
    "fmt"
    "github.com/jbitor/bittorrent"
    "github.com/jbitor/dht"
    "github.com/jbitor/webclient"
    "os"
    "log"
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

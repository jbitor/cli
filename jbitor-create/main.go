package main

import (
	"crypto/sha1"
	weakrand "math/rand"
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
			"%{color}%{level:4.4s} %{id:03x}%{color:reset} %{message}\n         %{longfunc}() in %{module}/%{shortfile}\n\n",
		)))
}

const PieceLength = 32768

func main() {
	if len(os.Args) != 2 {
		logger.Fatalf("Usage: %v PATH\n", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	path := os.Args[1]

	infoDict, err := bittorrent.GenerateTorrentMetaInfo(bittorrent.CreationOptions{
		Path:           path,
		PieceLength:    PieceLength,
		ForceMultiFile: false,
	})
	if err != nil {
		logger.Fatalf("Error generating torrent: %v\n", err)
		return
	}

	infoData, err := bencoding.Encode(infoDict)
	if err != nil {
		logger.Fatalf("Error encoding torrent infodict (for hashing): %v\n", err)
		return
	}

	torrentDict := bencoding.Dict{
		"info": infoDict,
		"nodes": bencoding.List{
			bencoding.List{
				bencoding.String("127.0.0.1"),
				bencoding.Int(6881),
			},
		},
	}

	torrentData, err := bencoding.Encode(torrentDict)

	if err != nil {
		logger.Fatalf("Error encoding torrent data: %v\n", err)
		return
	}

	hasher := sha1.New()
	hasher.Write(infoData)
	hash := hasher.Sum(nil)
	infoHash := bittorrent.BTID(hash)

	logger.Info("Generated torrent btih=%v.\n", infoHash)

	os.Stdout.Write(torrentData)
	os.Stdout.Sync()
}

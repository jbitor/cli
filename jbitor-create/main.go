package main
	
import (
    "crypto/sha1"
    "github.com/jbitor/bencoding"
    "github.com/jbitor/bittorrent"
    "os"
    "log"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stderr, "", 0)
}

const PieceLength = 32768

func main() {
    if len(os.Args) != 2 {
        logger.Fatalf("Usage: %v PATH\n", os.Args[0])
        return
    }

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

    logger.Printf("Generated torrent btih=%v.\n", infoHash)

    os.Stdout.Write(torrentData)
    os.Stdout.Sync()
}

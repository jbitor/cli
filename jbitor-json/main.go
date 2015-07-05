package main

import (
	"encoding/json"
	"io/ioutil"
	weakrand "math/rand"
	"os"
	"time"

	"github.com/jbitor/bencoding"
	"github.com/jbitor/cli/loggerconfig"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func init() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			"%{color}%{level:4.4s}%{color:reset} %{message}\n%{color}%{id:4.4x}%{color:reset} %{module} / %{shortfile} / %{longfunc}()\n",
		)))
}

func main() {
	loggerconfig.Use()

	if len(os.Args) == 1 {
		logger.Fatalf("Usage: %v from-bencoding|to-bencoding", os.Args[0])
		return
	}

	weakrand.Seed(time.Now().UTC().UnixNano())

	subcommand := os.Args[1]
	subcommandArgs := os.Args[2:]

	switch subcommand {
	case "from-bencoding":
		cmdJsonFromBencoding(subcommandArgs)
	case "to-bencoding":
		cmdJsonToBencoding(subcommandArgs)
	default:
		logger.Fatalf("Unknown torrent subcommand: %v", subcommand)
		return
	}

}

func cmdJsonFromBencoding(args []string) {
	if len(args) != 0 {
		logger.Fatalf("Usage: %v from-bencoding < FOO.torrent > FOO.bittorrent.json", os.Args[0])
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logger.Fatalf("Error reading stdin: %v", err)
		return
	}

	decoded, err := bencoding.Decode(data)
	if err != nil {
		logger.Fatalf("Error bdecoding stdin: %v", err)
		return
	}

	jsonable, err := decoded.ToJsonable()
	if err != nil {
		logger.Fatalf("Error converting bencoded value to jsonable: %v", err)
	}

	jsoned, err := json.Marshal(jsonable)
	if err != nil {
		logger.Fatalf("Error json-encoding data: %v", err)
		return
	}

	os.Stdout.Write(jsoned)
	os.Stdout.Sync()
}

func cmdJsonToBencoding(args []string) {
	if len(args) != 0 {
		logger.Fatalf("Usage: %v to-bencoding < FOO.bittorrent.json > FOO.torrent", os.Args[0])
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logger.Fatalf("Error reading stdin: %v", err)
		return
	}

	var decoded *interface{}
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		logger.Fatalf("Error decoding JSON from stdin: %v", err)
		return
	}

	bval, err := bencoding.FromJsonable(*decoded)
	if err != nil {
		logger.Fatalf("Error converting jsonable to bencodable: %v", err)
		return
	}

	encoded, err := bencoding.Encode(bval)
	if err != nil {
		logger.Fatalf("Error bencoding value: %v", err)
		return
	}

	os.Stdout.Write(encoded)
	os.Stdout.Sync()
}

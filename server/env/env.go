package env

import "flag"

type Flags struct {
	AudioDirectory string
	Address        string
	IndexName      string
}

var Options = Flags{}

func init() {
	audioDirectory := flag.String("audio", "audio", "the directory for storing audio content")
	address := flag.String("address", "localhost:8000", "the address to host on")
	indexName := flag.String("index-name", "index.pb", "the name for the index file")
	flag.Parse()

	Options = Flags{
		AudioDirectory: *audioDirectory,
		Address:        *address,
		IndexName:      *indexName,
	}
}

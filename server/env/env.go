package env

import (
	"strings"

	"github.com/LQR471814/music-player/server/logging"

	"github.com/hujun-open/sconf"
)

type Flags struct {
	AudioDirectory    string //the directory for storing audio content
	Address           string //the address to host on
	IndexName         string //the name for the index file
	Reset             bool   //reset all indexes
	PaletteResolution int    //the resolution (average between width and height) to store the color palette reference images with
}

var Options = Flags{}

func init() {
	defaultConfig := Flags{
		AudioDirectory:    "audio",
		Address:           ":6325",
		IndexName:         "index.pb",
		Reset:             false,
		PaletteResolution: 600,
	}

	conf, err := sconf.NewSConfCMDLine(defaultConfig, "config.yaml")
	if err != nil {
		logging.Error.Fatal(err)
	}
	ferr, aerr := conf.ReadwithCMDLine()
	if !strings.Contains(ferr.Error(), "config.yaml") && ferr != nil {
		logging.Warn.Println(ferr)
	}
	if aerr != nil {
		logging.Warn.Println(aerr)
	}

	Options = conf.GetConf()
}

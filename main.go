package main

import (
	"log"

	. "commonsyllabi/logger"
	"commonsyllabi/viewer"
)

func main() {
	InitLog()
	Log.Info().Msg("CC Viewer v0.1")

	var cc viewer.Cartridge
	cc = viewer.NewIMSCC()
	cc, err := cc.Load("test_files/canvas_large_1.3.imscc")
	if err != nil {
		log.Fatal(err)
	}
	Log.Info().Msg("successfully loaded cartridge")

	_, err = cc.AsObject()
	if err != nil {
		log.Fatal(err)
	}
}

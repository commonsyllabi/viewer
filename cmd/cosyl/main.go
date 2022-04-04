package main

import (
	"log"

	zero "commonsyllabi/internals/logger"
	"commonsyllabi/pkg/commoncartridge"
)

func main() {
	zero.InitLog()
	zero.Log.Info().Msg("CC Viewer v0.1")

	cc, err := commoncartridge.Load("test_files/canvas_large_1.3.imscc")
	if err != nil {
		log.Fatal(err) //-- todo remove and change to zerolog
	}
	zero.Log.Info().Msg("successfully loaded cartridge")

	_, err = cc.AsObject()
	if err != nil {
		log.Fatal(err)
	}
}

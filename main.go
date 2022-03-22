package main

import (
	"fmt"
	"log"

	"commonsyllabi/viewer"
)

func main() {
	fmt.Println("cc viewer 0.1")

	var cc viewer.Cartridge
	cc = viewer.NewIMSCC()
	cc, err := cc.Load("test_files/canvas_large_1.3.imscc")
	cc.Dump()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully loaded cartridge")
}

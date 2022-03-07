package main

import (
	"fmt"
	"log"

	"commonsyllabi/viewer"
)

func main() {
	fmt.Println("cc viewer 0.1")

	manifest, err := viewer.LoadFile("test_files/canvas_large_1.3.imscc")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("found manifest at: %s\n", manifest)

	err = viewer.ParseManifest(manifest)

	if err != nil {
		log.Fatal(err)
	}
}

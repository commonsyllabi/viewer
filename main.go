package main

import (
	"fmt"
	"log"

	"commonsyllabi/viewer"
)

func main() {
	fmt.Println("cc viewer 0.1")

	err := viewer.LoadFile("test_files/canvas_large_1.3.imscc")

	if err != nil {
		log.Fatal(err)
	}
}

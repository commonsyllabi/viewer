package main

import (
	"fmt"
	"log"

	"commonsyllabi/viewer"
)

func main() {
	fmt.Println("cc viewer 0.1")

	err := viewer.ReadFile()

	if err != nil {
		log.Fatal(err)
	}
}

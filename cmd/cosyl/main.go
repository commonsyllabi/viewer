package main

import (
	"flag"
	"fmt"
	"log"

	"commonsyllabi/pkg/commoncartridge"
)

var (
	debug    = flag.Bool("d", false, "debug output")
	metadata = flag.Bool("m", false, "show metadata")
	json     = flag.Bool("j", false, "dumps a serialized json representation")
)

func main() {
	flag.Parse()
	if *debug {
		fmt.Println("cosyl v0.1")
	}

	if flag.NArg() == 0 {
		log.Fatal("provide the path of the cartridge to be opened!")
	}

	inputFile := flag.Args()[0]

	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		fmt.Println("successfully loaded cartridge")
	}

	if *metadata {
		meta, err := cc.Metadata()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(meta)
	}

	if *json {
		obj, err := cc.AsObject()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(obj)
	}
}

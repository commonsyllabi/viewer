package main

import (
	"flag"
	"fmt"
	"log"

	"commonsyllabi/pkg/commoncartridge"
)

var (
	debug     = flag.Bool("d", false, "debug output")
	metadata  = flag.Bool("m", false, "shows metadata as serialized json")
	json      = flag.Bool("j", false, "dumps a serialized json representation")
	resources = flag.Bool("r", false, "lists all resources in the cartridge")
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

	if *resources {
		resources, err := cc.Resources()
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range resources {
			if r.Item.Identifierref == "" {
				r.Item.Title = "none"
			}
			if r.Resource.Href == "" {
				r.Resource.Href = "none"
			}

			fmt.Printf("type: %s\nfiles: %d \nhref: %s\nitem: %s\n\n", r.Resource.Type, len(r.Resource.File), r.Resource.Href, r.Item.Title)
		}
	}

	if *json {
		obj, err := cc.AsObject()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(obj)
	}
}

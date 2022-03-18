package specs

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type Manifest struct {
	XMLName       xml.Name      `xml:"manifest"`
	Metadata      Metadata      `xml:"metadata"`
	Organizations Organizations `xml:"organizations"`
	Resources     Resources     `xml:"resources"`
}

type Generic struct {
	XMLName xml.Name
}

func (m Manifest) TraverseItems(items []Item, root string) {
	for i, _ := range items {

		if items[i].Identifierref != "" {
			fmt.Printf("checking for idref: %v\n", items[i].Identifierref) //-- should i skip if there is no identifierref? probably
			//-- check if there are any related files in the top level directory

			files, err := ioutil.ReadDir(root)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				if strings.Contains(f.Name(), items[i].Identifierref) {
					fmt.Printf("idref %v - found matching file %s\n", items[i].Identifierref, f.Name())

					if f.IsDir() {
						fmt.Println("found dir, entering...")
						innerfiles, err := ioutil.ReadDir(filepath.Join(root, f.Name()))

						if err != nil {
							log.Fatal(err)
						}

						for _, inner := range innerfiles {

							if filepath.Ext(inner.Name()) == ".xml" {
								err = items[i].parseItem(filepath.Join(root, f.Name(), inner.Name()))

								if err != nil {
									log.Fatal(err)
								}
							} else {
								fmt.Printf("...appending %s\n", inner.Name())
								//-- note: we might want to store relative paths
								items[i].Filepaths = append(items[i].Filepaths, filepath.Join(root, f.Name(), inner.Name()))
							}

						}
					} else {
						items[i].parseItem(filepath.Join(root, f.Name()))

						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}

		if len(items[i].Items) > 0 {
			m.TraverseItems(items[i].Items, root)
		}
	}
}

func (m Manifest) prettyPrint() {
	s, _ := json.MarshalIndent(m, "", "\t")
	fmt.Printf("manifest struct: %s\n", s)

	fmt.Printf("title: %v\n", m.Metadata.LOM.General.Title.Value)
}

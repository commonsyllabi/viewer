package specs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func (m Manifest) TraverseItemOrg(itemOrg Item, root string) {
	m.traverseItemModules(itemOrg.Item, root)
}

func (m Manifest) traverseItemModules(itemModules []Item, root string) {
	for i, _ := range itemModules {
		m.traverseItems(itemModules[i].Item, root)
	}
}

func (m Manifest) traverseItems(items []Item, root string) {
	for i, _ := range items {
		for j, _ := range items[i].Item {
			if items[i].Item[j].Identifierref != "" {
				fmt.Printf("checking for idref: %v\n", items[i].Item[j].Identifierref) //-- should i skip if there is no identifierref? probably
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
									fmt.Printf("...skipping %s\n", inner.Name())
									//-- note: we might want to store relative paths
									// items[i].Filepaths = append(items[i].Filepaths, filepath.Join(root, f.Name(), inner.Name()))
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

			if len(items[i].Item) > 0 {
				m.traverseItems(items[i].Item, root)
			}
		}
	}
}

func (m Manifest) PrettyPrint() {
	s, _ := json.MarshalIndent(m, "", "\t")
	fmt.Printf("manifest struct: %s\n", s)

	fmt.Printf("title: %v\n", m.Metadata.Lom.General.Title)
}

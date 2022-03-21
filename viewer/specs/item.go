package specs

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Generic struct {
	XMLName xml.Name
}

func (i Item) parseItem(path string) error {

	//-- first open the file
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var gen Generic

	xml.Unmarshal(bytes, &gen)

	fmt.Printf("parsed: %v\n", gen.XMLName)

	//-- then big switch case with all the different possibilities
	//-- assignment
	//-- qti
	//-- topic
	//-- weblink
	//-- topicmeta

	return nil

}

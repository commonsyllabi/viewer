package specs

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Assignment struct {
	XMLName xml.Name
}

type Assessment struct {
	XMLName xml.Name
}

type Item struct {
	XMLName       xml.Name
	Identifier    string `xml:"identifier,attr"`
	Identifierref string `xml:"identifierref,attr"`
	Title         string `xml:"title"`
	Items         []Item `xml:"item"`
	Filepaths     []string
	Assignments   []Assignment
	Assessments   []Assessment
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

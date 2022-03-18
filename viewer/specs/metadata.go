package specs

import "encoding/xml"

type Title struct {
	XMLName xml.Name `xml:"title"`
	Value   string   `xml:"string"`
}

type General struct {
	XMLName xml.Name `xml:"general"`
	Title   Title    `xml:"title"`
}

type LOM struct {
	XMLName xml.Name `xml:"lom"`
	General General  `xml:"general"`
}

type Metadata struct {
	XMLName       xml.Name `xml:"metadata"`
	Schema        string   `xml:"schema"`
	Schemaversion string   `xml:"schemaversion"`
	LOM           LOM      `xml:"lom"`
}

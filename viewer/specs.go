package viewer

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

type Manifest struct {
	XMLName       xml.Name      `xml:"manifest"`
	Metadata      Metadata      `xml:"metadata"`
	Organizations Organizations `xml:"organizations"`
	// Resources     Resources     `xml:"resources"`
}

type Organizations struct {
	XMLName       xml.Name       `xml:"organizations"`
	Organizations []Organization `xml:"organization"`
}

type Organization struct {
	XMLName    xml.Name `xml:"organization"` //-- each of these is a learning module
	Identifier string   `xml:"identifier,attr"`
	Structure  string   `xml:"structure,attr"`
	Items      []Item   `xml:"item"`
}

type Item struct {
	XMLName       xml.Name
	Identifier    string `xml:"identifier,attr"`
	Identifierref string `xml:"identifierref,attr"`
	Title         string `xml:"title"` //-- problem with circular declaration of item within item within...
	Items         []Item `xml:"item"`
}

type Resources struct {
	XMLName   xml.Name   `xml:"resources"`
	Resources []Resource `xml:"resource"`
}

type Resource struct {
	XMLName     xml.Name `xml:"resource"`
	Identifier  string   `xml:"identifier,attr"`
	Type        string   `xml:"type,attr"`
	IntendedUse string   `xml:"intendeduse,attr"`
	Href        string   `xml:"href,attr"`
	File        File     `xml:"file"`
}

type File struct {
	XMLName xml.Name `xml:"file"`
	Href    string   `xml:"href,attr"`
}

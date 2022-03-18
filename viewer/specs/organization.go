package specs

import "encoding/xml"

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

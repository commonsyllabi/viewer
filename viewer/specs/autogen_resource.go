// Code generated by zek; DO NOT EDIT.

package specs

import "encoding/xml"

// Resource was generated 2022-03-21 14:18:27 by pierre on pierre-mac.
type Resource struct {
	XMLName     xml.Name `xml:"resource"`
	Text        string   `xml:",chardata"`
	Identifier  string   `xml:"identifier,attr"`
	Type        string   `xml:"type,attr"`
	Href        string   `xml:"href,attr"`
	Intendeduse string   `xml:"intendeduse,attr"`
	File        []struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"file"`
	Variant struct {
		Text          string `xml:",chardata"`
		Identifier    string `xml:"identifier,attr"`
		Identifierref string `xml:"identifierref,attr"`
		Metadata      string `xml:"metadata"`
	} `xml:"variant"`
	Dependency struct {
		Text          string `xml:",chardata"`
		Identifierref string `xml:"identifierref,attr"`
	} `xml:"dependency"`
}

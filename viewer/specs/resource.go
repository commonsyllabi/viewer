package specs

import "encoding/xml"

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

//-- heads up, when parsing assessments, QTI and Meta are separate but they share a common id (id of meta is folder of qti)
//-- can probably be combined into 1

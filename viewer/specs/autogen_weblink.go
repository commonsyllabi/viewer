// Code generated by zek; DO NOT EDIT.

package specs

import "encoding/xml"

// WebLink was generated 2022-03-29 17:51:56 by pierre on archpierre.
type WebLink struct {
	XMLName        xml.Name `xml:"webLink"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Title          string   `xml:"title"`
	URL            struct {
		Text           string `xml:",chardata"`
		Href           string `xml:"href,attr"`
		Target         string `xml:"target,attr"`
		WindowFeatures string `xml:"windowFeatures,attr"`
	} `xml:"url"`
}

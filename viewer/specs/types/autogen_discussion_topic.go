// Code generated by zek; DO NOT EDIT.

package specs

import "encoding/xml"

// Topic was generated 2022-03-21 16:49:23 by pierre on pierre-mac.
type Topic struct {
	XMLName        xml.Name `xml:"topic"`
	Chardata       string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Title          string   `xml:"title"`
	Text           struct {
		Text     string `xml:",chardata"`
		Texttype string `xml:"texttype,attr"`
	} `xml:"text"`
	Attachments struct {
		Text       string `xml:",chardata"`
		Attachment []struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
		} `xml:"attachment"`
	} `xml:"attachments"`
}

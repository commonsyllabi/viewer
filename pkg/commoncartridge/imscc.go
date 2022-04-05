package commoncartridge

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type IMSCC struct {
	Reader zip.Reader
	Path   string
}

//-- Load creates a cartridge from a given path
func Load(path string) (IMSCC, error) {
	cc := IMSCC{}

	r, err := zip.OpenReader(path)
	if err != nil {
		return cc, err
	}

	cc.Reader = r.Reader
	cc.Path = path

	return cc, nil
}

type Metadata struct {
	Title                string
	Schema               string
	SchemaVersion        string
	Language             string
	Description          string
	Keyword              string
	Date                 string
	Copyright            string
	CopyrightDescription string
}

func (cc IMSCC) Metadata() (string, error) {

	m, err := cc.ParseManifest()

	if err != nil {
		return "", err
	}

	meta := Metadata{
		m.Metadata.Lom.General.Title.String.Text,
		m.Metadata.Schema,
		m.Metadata.Schemaversion,
		m.Metadata.Lom.General.Language,
		m.Metadata.Lom.General.Description.String.Text,
		m.Metadata.Lom.General.Keyword.String.Text,
		m.Metadata.Lom.LifeCycle.Contribute.Date.DateTime,
		m.Metadata.Lom.Rights.CopyrightAndOtherRestrictions.Value,
		m.Metadata.Lom.Rights.Description.String,
	}

	serialized, err := json.Marshal(meta)

	if err != nil {
		return string(serialized), nil
	}

	return string(serialized), nil
}

func (cc IMSCC) Title() string {
	title := "--undefined--"

	var m Manifest
	m, err := cc.ParseManifest()

	if err != nil {
		fmt.Printf("Error parsing Manifest: %v\n", err)
	}

	title = m.Metadata.Lom.General.Title.String.Text
	return title
}

func (cc IMSCC) Dump() []string {
	dump := make([]string, len(cc.Reader.File))
	for _, f := range cc.Reader.File {
		dump = append(dump, f.Name)
	}
	return dump
}

func (cc IMSCC) ParseManifest() (Manifest, error) {

	var manifest Manifest
	var path string
	for _, f := range cc.Reader.File {
		if strings.Contains(f.Name, "imsmanifest.xml") {
			path = f.Name
		}
	}

	file, err := cc.Reader.Open(path)

	if err != nil {
		fmt.Printf("Error in opening manifest: %v\n", cc.Path)
		return manifest, err
	}

	bytesArray, err := io.ReadAll(file)
	if err != nil {
		return manifest, err
	}

	xml.Unmarshal(bytesArray, &manifest)

	manifest.ResolveItems()

	return manifest, nil
}

func (cc IMSCC) AsObject() ([]byte, error) {
	var obj []byte

	m, err := cc.ParseManifest()
	if err != nil {
		return obj, err
	}

	obj, err = json.Marshal(m)
	if err != nil {
		return obj, err
	}
	return obj, err
}

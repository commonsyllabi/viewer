package commoncartridge

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
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

func (cc IMSCC) ParseManifest() (Manifest, error) {

	var manifest Manifest
	var path string
	for _, f := range cc.Reader.File {
		if strings.Contains(f.Name, "imsmanifest.xml") {
			path = f.Name
			break
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

type FullResource struct {
	Resource Resource
	Item     Item
}

func (cc IMSCC) Resources() ([]FullResource, error) {
	resources := make([]FullResource, 0)

	m, err := cc.ParseManifest()

	if err != nil {
		return resources, err
	}

	for _, r := range m.Resources.Resource {
		res := FullResource{}
		res.Resource = r

		item, err := m.FindItem(r.Identifier)
		if err != nil {
			return resources, err
		}
		res.Item = item
		resources = append(resources, res)
	}

	return resources, nil
}

func (cc IMSCC) Weblinks() ([]WebLink, error) {
	weblinks := make([]WebLink, 0)

	m, err := cc.ParseManifest()

	if err != nil {
		return weblinks, err
	}

	re, err := regexp.Compile(`imswl_xmlv1p\d`)
	if err != nil {
		return weblinks, err
	}

	for _, r := range m.Resources.Resource {
		match := re.Find([]byte(r.Type))

		if match != nil {

			var wl WebLink
			var path string
			for _, f := range cc.Reader.File {
				if strings.Contains(f.Name, r.Identifier+".xml") {
					path = f.Name
					break
				}
			}

			file, err := cc.Reader.Open(path)
			if err != nil {
				return weblinks, err
			}

			bytesArray, err := io.ReadAll(file)
			if err != nil {
				return weblinks, err
			}

			xml.Unmarshal(bytesArray, &wl)

			weblinks = append(weblinks, wl)
		}
	}

	return weblinks, nil
}

func (cc IMSCC) Dump() []string {
	dump := make([]string, len(cc.Reader.File))
	for _, f := range cc.Reader.File {
		dump = append(dump, f.Name)
	}
	return dump
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

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

func (cc IMSCC) Assignments() ([]Assignment, error) {
	assignments := make([]Assignment, 0)

	paths, err := cc.findResourcesByType(`assignment_xmlv1p\d`)

	if err != nil {
		return assignments, err
	}

	for _, p := range paths {
		file, err := cc.Reader.Open(p)
		if err != nil {
			return assignments, err
		}

		bytesArray, err := io.ReadAll(file)
		if err != nil {
			return assignments, err
		}

		var a Assignment
		xml.Unmarshal(bytesArray, &a)
		//-- necessary check for avoiding other files in the folder that are returned by the findResourcesByType() (e.g. `assignment.xml` also has `assignment_meta.html`)
		if a.XMLName.Local == "assignment" {
			assignments = append(assignments, a)
		}
	}

	return assignments, nil
}

func (cc IMSCC) LTIs() ([]CartridgeBasicltiLink, error) {
	qtis := make([]CartridgeBasicltiLink, 0)

	paths, err := cc.findResourcesByType(`imsbasiclti_xmlv1p\d`)

	if err != nil {
		return qtis, err
	}

	for _, p := range paths {
		file, err := cc.Reader.Open(p)
		if err != nil {
			return qtis, err
		}

		bytesArray, err := io.ReadAll(file)
		if err != nil {
			return qtis, err
		}

		var qti CartridgeBasicltiLink
		xml.Unmarshal(bytesArray, &qti)
		qtis = append(qtis, qti)
	}

	return qtis, nil
}

func (cc IMSCC) QTIs() ([]Questestinterop, error) {
	qtis := make([]Questestinterop, 0)

	paths, err := cc.findResourcesByType(`imsqti_xmlv1p\d`)

	if err != nil {
		return qtis, err
	}

	for _, p := range paths {
		file, err := cc.Reader.Open(p)
		if err != nil {
			return qtis, err
		}

		bytesArray, err := io.ReadAll(file)
		if err != nil {
			return qtis, err
		}

		var qti Questestinterop
		xml.Unmarshal(bytesArray, &qti)

		if qti.XMLName.Local == "questestinterop" {
			qtis = append(qtis, qti)
		}
	}

	return qtis, nil
}

func (cc IMSCC) Topics() ([]Topic, error) {
	topics := make([]Topic, 0)

	paths, err := cc.findResourcesByType(`imsdt_xmlv1p\d`)

	if err != nil {
		return topics, err
	}

	for _, p := range paths {
		file, err := cc.Reader.Open(p)
		if err != nil {
			return topics, err
		}

		bytesArray, err := io.ReadAll(file)
		if err != nil {
			return topics, err
		}

		var t Topic
		xml.Unmarshal(bytesArray, &t)
		topics = append(topics, t)
	}

	return topics, nil
}

func (cc IMSCC) Weblinks() ([]WebLink, error) {
	weblinks := make([]WebLink, 0)

	paths, err := cc.findResourcesByType(`imswl_xmlv1p\d`)

	if err != nil {
		return weblinks, err
	}

	for _, p := range paths {
		file, err := cc.Reader.Open(p)
		if err != nil {
			return weblinks, err
		}

		bytesArray, err := io.ReadAll(file)
		if err != nil {
			return weblinks, err
		}

		var wl WebLink
		xml.Unmarshal(bytesArray, &wl)
		weblinks = append(weblinks, wl)
	}

	return weblinks, nil
}

func (cc IMSCC) Find(id string) (Resource, error) {
	m, err := cc.ParseManifest()

	if err != nil {
		return Resource{}, err
	}

	for _, r := range m.Resources.Resource {
		if r.Identifier == id {
			return r, nil
		}
	}

	return Resource{}, fmt.Errorf("could not find resource with id: %v", id)
}

// FindFile takes an ID and returns the zip.File from the cartridge's Reader
func (cc IMSCC) FindFile(id string) (zip.File, error) {
	file := zip.File{}
	m, err := cc.ParseManifest()

	if err != nil {
		return file, err
	}

	for _, r := range m.Resources.Resource {
		if r.Type == "webcontent" && r.Identifier == id {
			for _, f := range cc.Reader.File {
				if strings.Contains(f.Name, r.Href) {
					if err != nil {
						return *f, err
					}
					return *f, nil
				}
			}
		}
	}

	return file, nil
}

// findResourcesByType takes a regex pattern and returns a slice of paths of files who match the pattern
func (cc IMSCC) findResourcesByType(pattern string) ([]string, error) {
	paths := make([]string, 0)

	re, err := regexp.Compile(pattern)
	if err != nil {
		return paths, err
	}

	m, err := cc.ParseManifest()

	if err != nil {
		return paths, err
	}

	for _, r := range m.Resources.Resource {
		match := re.Find([]byte(r.Type))

		if match != nil {
			for _, f := range cc.Reader.File {
				if !f.FileInfo().IsDir() && strings.Contains(f.Name, r.Identifier) {
					paths = append(paths, f.Name)
				}
			}
		}
	}

	return paths, nil
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

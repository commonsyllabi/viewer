// CommonCartridge allows you to manipulate IMSCC-compliant Common Cartrdiges through its manifest, resources and associated files.
package commoncartridge

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/commonsyllabi/viewer/pkg/commoncartridge/types"
)

// IMSCC loads the IMSCC-specific cartridge into a zip.Reader from the given Path
type IMSCC struct {
	Reader zip.Reader
	Path   string
}

// Load returns a cartridge created from a given path and holds it into a zip.Reader
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

// Parse Manifest finds and marshals the imsccmanifest.xml file into the Manifest struct
func (cc IMSCC) ParseManifest() (types.Manifest, error) {

	var manifest types.Manifest
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

// Title returns the name of the cartridge
func (cc IMSCC) Title() string {
	title := "--undefined--"

	var m types.Manifest
	m, err := cc.ParseManifest()

	if err != nil {
		fmt.Printf("Error parsing Manifest: %v\n", err)
	}

	title = m.Metadata.Lom.General.Title.String.Text
	return title
}

//-- TODO use an autogenerated version of this struct?
// Metadata is a user-friendly representation of the Metadata field of the Manifest
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

// Metadata returns a user-friendly, stringified, JSON-encoded version of the Metadata field.
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

// FullResource is a union of a Resource and the Item that refers to it
type FullResource struct {
	Resource interface{}
	Item     types.Item //-- as is, this might not be very meaningful; might need to store every other parent Item
}

// Resources returns a slice of all FullResources, each containing a resource and the item it can belong to
func (cc IMSCC) Resources() ([]FullResource, error) {
	resources := make([]FullResource, 0)

	m, err := cc.ParseManifest()

	if err != nil {
		return resources, err
	}

	for _, r := range m.Resources.Resource {
		res := FullResource{}
		res.Resource, err = cc.Find(r.Identifier) //todo: resolve the actual resource in there
		if err != nil {
			return resources, err
		}

		item, err := m.FindItem(r.Identifier)
		if err != nil {
			return resources, err
		}
		res.Item = item
		resources = append(resources, res)
	}

	return resources, nil
}

// Assignments returns a slice of all resources of type assignment_xmlv1p\d
func (cc IMSCC) Assignments() ([]types.Assignment, error) {
	assignments := make([]types.Assignment, 0)

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

		var a types.Assignment
		xml.Unmarshal(bytesArray, &a)
		//-- necessary check for avoiding other files in the folder that are returned by the findResourcesByType() (e.g. `assignment.xml` also has `assignment_meta.html`)
		if a.XMLName.Local == "assignment" {
			assignments = append(assignments, a)
		}
	}

	return assignments, nil
}

// LTIs returns a slice of all resources of type imsbasiclti_xmlv1p\d
func (cc IMSCC) LTIs() ([]types.CartridgeBasicltiLink, error) {
	qtis := make([]types.CartridgeBasicltiLink, 0)

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

		var qti types.CartridgeBasicltiLink
		xml.Unmarshal(bytesArray, &qti)
		qtis = append(qtis, qti)
	}

	return qtis, nil
}

// QTIs returns a slice of all resources of type imsqti_xmlv1p\d
func (cc IMSCC) QTIs() ([]types.Questestinterop, error) {
	qtis := make([]types.Questestinterop, 0)

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

		var qti types.Questestinterop
		xml.Unmarshal(bytesArray, &qti)

		if qti.XMLName.Local == "questestinterop" {
			qtis = append(qtis, qti)
		}
	}

	return qtis, nil
}

// Topics returns a slice of all resources of type imsdt_xmlv1p\d
func (cc IMSCC) Topics() ([]types.Topic, error) {
	topics := make([]types.Topic, 0)

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

		var t types.Topic
		xml.Unmarshal(bytesArray, &t)
		topics = append(topics, t)
	}

	return topics, nil
}

// Weblnks returns a slice of all resources of type imswl_xmlv1p\d
func (cc IMSCC) Weblinks() ([]types.WebLink, error) {
	weblinks := make([]types.WebLink, 0)

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

		var wl types.WebLink
		xml.Unmarshal(bytesArray, &wl)
		weblinks = append(weblinks, wl)
	}

	return weblinks, nil
}

// Find takes an id and returns the resource associated with it
func (cc IMSCC) Find(id string) (interface{}, error) {
	m, err := cc.ParseManifest()

	if err != nil {
		return types.Resource{}, err
	}

	//-- find the type, then marshal into the appropriate struct
	//-- otherwise return the resource
	for _, r := range m.Resources.Resource {

		if r.Identifier == id {

			var path string
			if r.Href != "" {
				path = r.Href
			} else { //-- find the href in the first File nodes of XML type
				for _, f := range r.File {
					if strings.Contains(f.Href, ".xml") {
						path = f.Href
						break
					}
				}
			}

			// todo decide to handle _fallback here or in client
			// otherwise, return the resource as is
			if path == "" {
				return r, nil
			}

			file, err := cc.Reader.Open(path)
			if err != nil {
				return r, err
			}

			bytes, err := io.ReadAll(file)
			if err != nil {
				return r, err
			}

			switch r.Type {
			case "imsdt_xmlv1p0", "imsdt_xmlv1p1", "imsdt_xmlv1p2", "imsdt_xmlv1p3":
				// fmt.Printf("found topic %v", resource.File)
				var t types.Topic
				err = xml.Unmarshal(bytes, &t)
				if err != nil {
					return t, nil
				}
				return t, nil

			case "webcontent":
				// fmt.Printf("found webcontent %v", resource.File)
				return r, nil
			case "imswl_xmlv1p0", "imswl_xmlv1p1", "imswl_xmlv1p2", "imswl_xmlv1p3":
				// fmt.Printf("found weblink %v", resource.File)
				var wl types.WebLink
				err = xml.Unmarshal(bytes, &wl)
				if err != nil {
					return wl, nil
				}
				return wl, nil

			case "assignment_xmlv1p0", "assignment_xmlv1p1", "assignment_xmlv1p2", "assignment_xmlv1p3":
				// fmt.Printf("found assignment %v", resource.File)
				var a types.Assignment
				err = xml.Unmarshal(bytes, &a)
				if err != nil {
					return a, nil
				}
				return a, nil

			case "imsqti_xmlv1p2/imscc_xmlv1p1/assessment", "imsqti_xmlv1p2/imscc_xmlv1p2/assessment",
				"imsqti_xmlv1p2/imscc_xmlv1p3/assessment":
				// fmt.Printf("found question bank %v", resource.File)
				var qti types.Questestinterop
				err = xml.Unmarshal(bytes, &qti)
				if err != nil {
					return qti, nil
				}
				return qti, nil
			case "imsbasiclti_xmlv1p0", "imsbasiclti_xmlv1p1", "imsbasiclti_xmlv1p2":
				//-- fmt.Printf("found LTI %v\n, resource.File")
				var lti types.CartridgeBasicltiLink
				err = xml.Unmarshal(bytes, &lti)
				if err != nil {
					return lti, nil
				}
				return lti, nil
			case "associatedcontent/imscc_xmlv1p0/learning-application-resource", "associatedcontent/imscc_xmlv1p1/learning-application-resource", "associatedcontent/imscc_xmlv1p2/learning-application-resource",
				"associatedcontent/imscc_xmlv1p3/learning-application-resource":
				return r, nil
			default:
				return r, fmt.Errorf("no matching type found: %s", r.Type)
			}

		}
	}

	return types.Resource{}, fmt.Errorf("could not find resource with id: %v", id)
}

// FindFile takes an ID and returns the corresponding file as a byte slice
func (cc IMSCC) FindFile(id string) ([]byte, error) {
	var file bytes.Buffer
	m, err := cc.ParseManifest()

	if err != nil {
		return file.Bytes(), err
	}

	for _, r := range m.Resources.Resource {
		if r.Identifier == id {
			//-- directly go through the child []File and read from the href there

			b, err := cc.Reader.Open(r.File[0].Href)
			if err != nil {
				return file.Bytes(), err
			}

			bytes, err := io.ReadAll(b)
			if err != nil {
				return file.Bytes(), err
			}

			_, err = file.Write(bytes)
			if err != nil {
				return file.Bytes(), err
			}

			return file.Bytes(), nil
		}
	}

	return file.Bytes(), fmt.Errorf("couldn't find file for id %s", id)
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

// Dump returns the string representation of the Manifest
func (cc IMSCC) Dump() []string {
	dump := make([]string, len(cc.Reader.File))
	for _, f := range cc.Reader.File {
		dump = append(dump, f.Name)
	}
	return dump
}

// AsObject returns the JSON-encoded string representation of the Manifest
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

// CommonCartridge allows you to manipulate IMSCC-compliant Common Cartrdiges through its manifest, resources and associated files.
package commoncartridge

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"strings"

	"github.com/commonsyllabi/viewer/pkg/commoncartridge/types"
)

// IMSCC loads the IMSCC-specific cartridge into a zip.Reader from the given Path
type IMSCC struct {
	Reader   zip.Reader
	Path     string
	manifest types.Manifest
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
	cc.manifest, err = cc.parseManifest()

	return cc, err
}

func (cc IMSCC) Manifest() (types.Manifest, error) {
	return cc.manifest, nil
}

// Title returns the name of the cartridge
func (cc IMSCC) Title() string {
	return cc.manifest.Metadata.Lom.General.Title.String.Text
}

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

	meta := Metadata{
		cc.manifest.Metadata.Lom.General.Title.String.Text,
		cc.manifest.Metadata.Schema,
		cc.manifest.Metadata.Schemaversion,
		cc.manifest.Metadata.Lom.General.Language,
		cc.manifest.Metadata.Lom.General.Description.String.Text,
		cc.manifest.Metadata.Lom.General.Keyword.String.Text,
		cc.manifest.Metadata.Lom.LifeCycle.Contribute.Date.DateTime,
		cc.manifest.Metadata.Lom.Rights.CopyrightAndOtherRestrictions.Value,
		cc.manifest.Metadata.Lom.Rights.Description.String,
	}

	serialized, err := json.Marshal(meta)

	if err != nil {
		return string(serialized), err
	}

	return string(serialized), nil
}

// FullItem is a union of an Item and all resources that refer to it
type FullItem struct {
	Resources []types.Resource
	Item      types.Item
	Children  []FullItem
}

// Items returns all items with their associated resources. It goes through each item at a  level n and looks for full items at the level n-1
func (cc IMSCC) Items() ([]FullItem, error) {
	items := make([]FullItem, 0)

	//-- A CC always have only one top level item, so we can directly jump to its children
	for _, i := range cc.manifest.Organizations.Organization.Item.Item {

		full, err := cc.traverseItems(i.Item)

		if err != nil {
			return items, err
		}

		items = append(items, full...)
	}

	return items, nil
}

func (cc IMSCC) traverseItems(items []types.Item) ([]FullItem, error) {
	full := make([]FullItem, 0)

	// -- go through all children
	for _, i := range items {
		var f FullItem
		f.Item = i

		//-- add all resources
		for _, r := range cc.manifest.Resources.Resource {
			if strings.Contains(r.Identifier, i.Identifierref) {
				f.Resources = append(f.Resources, r)
			}
		}

		//-- if it has children, go through
		if len(i.Item) > 0 {
			children, err := cc.traverseItems(i.Item)

			f.Children = append(f.Children, children...)

			if err != nil {
				return full, nil
			}
		}

		full = append(full, f)

	}

	return full, nil
}

// FullResource is a union of a Resource and the Item that refers to it
type FullResource struct {
	Resource interface{}
	Item     types.Item
}

// Resources returns a slice of all FullResources, each containing a resource and the item it can belong to
func (cc IMSCC) Resources() ([]FullResource, error) {
	resources := make([]FullResource, 0)

	for _, r := range cc.manifest.Resources.Resource {
		res := FullResource{}
		found, err := cc.Find(r.Identifier)
		res.Resource = found

		if err != nil {
			return resources, err
		}

		item, err := cc.FindItem(r.Identifier)
		if err != nil {
			return resources, err
		}
		res.Item = item
		resources = append(resources, res)
	}

	return resources, nil
}

// FindItem returns that the item that the resources points to, or returns nil
func (cc *IMSCC) FindItem(id string) (types.Item, error) {

	var item types.Item
	var err error
	for _, i := range cc.manifest.Organizations.Organization.Item.Item {
		item, err = findItem(i.Item, id)

		if err != nil {
			return item, nil
		}

		if item.Identifierref == id {
			return item, nil
		}
	}

	return item, err
}

func findItem(items []types.Item, id string) (types.Item, error) {
	// fmt.Printf("looking for %s in %d items\n", id, len(items))
	var item types.Item
	var err error
	for i := range items {
		if items[i].Identifierref == id {
			return items[i], nil
		}

		if len(items[i].Item) > 0 {
			item, err = findItem(items[i].Item, id)
		}
	}

	return item, err
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

	//-- find the type, then marshal into the appropriate struct
	//-- otherwise return the resource
	for _, r := range cc.manifest.Resources.Resource {

		if r.Identifier == id {

			var path string
			if r.Href != "" {
				path = r.Href
			} else { //-- find the href in the first File nodes of XML type
				for _, f := range r.File {
					path = f.Href
				}
			}

			// todo should `_fallback` resource be appended to parent resource?
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
func (cc IMSCC) FindFile(id string) (fs.File, error) {
	var file fs.File

	for _, r := range cc.manifest.Resources.Resource {
		if r.Identifier == id {
			//-- directly go through the child []File and read from the href there

			f, err := cc.Reader.Open(r.File[0].Href)
			if err != nil {
				return f, err
			}

			return f, nil
		}
	}

	return file, fmt.Errorf("couldn't find file for id %s", id)
}

// findResourcesByType takes a regex pattern and returns a slice of paths of files who match the pattern
func (cc IMSCC) findResourcesByType(pattern string) ([]string, error) {
	paths := make([]string, 0)

	re, err := regexp.Compile(pattern)
	if err != nil {
		return paths, err
	}

	for _, r := range cc.manifest.Resources.Resource {
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

// parseManifest finds and marshals the imsccmanifest.xml file into the Manifest struct
func (cc IMSCC) parseManifest() (types.Manifest, error) {

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

	return manifest, nil
}

// Dump returns the string representation of the Manifest
func (cc IMSCC) Dump() []string {
	dump := make([]string, len(cc.Reader.File))
	for _, f := range cc.Reader.File {
		dump = append(dump, f.Name)
	}
	return dump
}

// MarshalJSON returns the JSON-encoded string representation of the Manifest
func (cc IMSCC) MarshalJSON() ([]byte, error) {
	var obj []byte

	obj, err := json.Marshal(cc.manifest)
	if err != nil {
		return obj, err
	}
	return obj, err
}

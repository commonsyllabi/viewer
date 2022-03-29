package viewer

import (
	"archive/zip"
	. "commonsyllabi/logger"
	"commonsyllabi/viewer/specs"
	"encoding/json"
	"encoding/xml"
	"io"
)

type IMSCC struct {
	Reader zip.Reader
	Path   string
}

func NewIMSCC() IMSCC {
	Log.Debug().Msg("Creating new IMSCC")
	var ims IMSCC
	return ims
}

//-- given a particular path, assigns a reader to a Cartridge
//-- and returns it
func (cc IMSCC) Load(path string) (Cartridge, error) {
	Log.Debug().Msg("Decompressing files")

	r, err := zip.OpenReader(path)
	if err != nil {
		return cc, err
	}

	cc.Reader = r.Reader
	cc.Path = path

	return cc, nil
}

func (cc IMSCC) Title() string {
	title := "--undefined--"

	var m specs.Manifest
	m, err := cc.ParseManifest()

	if err != nil {
		Log.Error().Msg("Error parsing Manifest")
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

//-- pointer receivers (*IMSCC) can modify the struct instance,
//-- while the value receivers can't change it
func (cc IMSCC) ParseManifest() (specs.Manifest, error) {

	var manifest specs.Manifest
	file, err := cc.Reader.Open("imsmanifest.xml")

	if err != nil {
		Log.Debug().Str("Error in opening manifest.", cc.Path)
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

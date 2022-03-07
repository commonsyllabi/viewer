package viewer

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Title struct {
	XMLName xml.Name `xml:"title"`
	Value   string   `xml:"string"`
}

type General struct {
	XMLName xml.Name `xml:"general"`
	Title   Title    `xml:"title"`
}

type LOM struct {
	XMLName xml.Name `xml:"lom"`
	General General  `xml:"general"`
}

type Metadata struct {
	XMLName       xml.Name `xml:"metadata"`
	Schema        string   `xml:"schema"`
	Schemaversion string   `xml:"schemaversion"`
	LOM           LOM      `xml:"lom"`
}

type Manifest struct {
	XMLName       xml.Name      `xml:"manifest"`
	Metadata      Metadata      `xml:"metadata"`
	Organizations Organizations `xml:"organizations"`
	// Resources     Resources     `xml:"resources"`
}

type Organizations struct {
	XMLName       xml.Name       `xml:"organizations"`
	Organizations []Organization `xml:"organization"`
}

type Organization struct {
	XMLName    xml.Name `xml:"organization"` //-- each of these is a learning module
	Identifier string   `xml:"identifier,attr"`
	Structure  string   `xml:"structure,attr"`
	Items      []Item   `xml:"item"`
}

type Item struct {
	XMLName       xml.Name `xml:"item"`
	Identifier    string   `xml:"identifier,attr"`
	Identifierref string   `xml:"identifierref,attr"`
	Title         string   `xml:"title"` //-- problem with circular declaration of item within item within...
}

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

//-- from https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func Unzip(src, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	filename := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))
	dest = filepath.Join(dest, filename)
	os.MkdirAll(dest, os.ModePerm)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()

		if err != nil {
			return err
		}

		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal) vulnerability
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)

			if err != nil {
				return err
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

			if err != nil {
				return err
			}

			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	file, err := os.OpenFile(filepath.Join(dest, "imsmanifest.xml"), os.O_RDONLY, os.ModePerm)
	if errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	defer file.Close()

	return file.Name(), nil
}

func LoadFile(path string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	dest := filepath.Join(pwd, "tmp")
	manifest, err := Unzip(path, dest)

	if err != nil {
		return "", err
	}

	return manifest, nil
}

func ParseManifest(manifest string) error {
	file, err := os.OpenFile(manifest, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Printf("opened manifest %s\n", file.Name())

	bytesArray, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var manif Manifest

	xml.Unmarshal(bytesArray, &manif)

	s, _ := json.MarshalIndent(manif, "", "\t")
	fmt.Printf("manifest struct: %s\n", s)

	fmt.Printf("title: %v\n", manif.Metadata.LOM.General.Title.Value)
	return nil
}

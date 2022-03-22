package viewer

import (
	"archive/zip"
	"commonsyllabi/viewer/specs"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type IMSCC struct {
	Reader zip.Reader
}

func NewIMSCC() IMSCC {
	var ims IMSCC
	return ims
}

//-- given a particular path, assigns a reader to a Cartridge
//-- and returns it
func (cc IMSCC) Load(path string) (Cartridge, error) {
	fmt.Println("decompressing files")

	r, err := zip.OpenReader(path)
	if err != nil {
		return cc, err
	}

	cc.Reader = r.Reader

	return cc, nil
}

func (cc IMSCC) Dump() {
	for _, f := range cc.Reader.File {
		fmt.Println(f.Name)
	}
}

//-- pointer receivers (*IMSCC) can modify the struct instance,
//-- while the value receivers can't change it
func (cc IMSCC) ParseManifest(root string) error {
	//-- todo, here, directy open from the zip reader
	file, err := os.OpenFile(filepath.Join(root, "imsmanifest.xml"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Printf("opened manifest %s\n", file.Name())

	bytesArray, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var manifest specs.Manifest

	xml.Unmarshal(bytesArray, &manifest)

	manifest.ResolveItems()
	manifest.PrettyPrint()

	return nil
}

func (cc IMSCC) AsObject() []byte {
	var obj []byte
	return obj
}

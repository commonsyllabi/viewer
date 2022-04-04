package commoncartridge

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const imscc_path = "./test_files/canvas_large_1.3.imscc"

// have setup() to create tmp files, with e.g. different schema versions
// check go test tables

func TestLoadEmpty(t *testing.T) {
	_, err := Load("")

	if err == nil {
		t.Fail()
	}
}

func TestParseManifest(t *testing.T) {
	// declaring cc as Cartridge to test that the return value implements Cartridge interface
	var cc Cartridge = load(t, imscc_path)

	manifest, err := cc.ParseManifest()

	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(manifest).Kind() != reflect.Struct {
		t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(manifest).Kind())
	}
}

func TestLoadCorrect(t *testing.T) {
	cc := load(t, imscc_path)

	var empty IMSCC
	if reflect.DeepEqual(cc, empty) {
		t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
	}
}

func TestLoadAll(t *testing.T) {
	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(filepath.Join(cwd, "./test_files"))

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Test loading %d cartridges\n", len(files))
	var i int = 0
	for _, file := range files {
		i++
		if file.IsDir() {
			continue
		}

		cc, err := Load(filepath.Join("./test_files", file.Name()))

		if err != nil {
			t.Error(err)
		}

		var empty IMSCC
		if reflect.DeepEqual(cc, empty) {
			t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
		}

		if cc.Title() == "" {
			t.Error("Cartridge Title should not be empty!")
		}

		fmt.Printf("Parsed %d/%d - %s\n", i, len(files), cc.Title())
	}
}

func TestParseManifestAll(t *testing.T) {
	cwd, _ := os.Getwd()
	files, err := ioutil.ReadDir(filepath.Join(cwd, "./test_files"))

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("testing %d cartridges\n", len(files))
	var i int = 0
	for _, file := range files {
		i++
		if file.IsDir() {
			continue
		}

		cc, err := Load(filepath.Join("./test_files", file.Name()))

		if err != nil {
			t.Error(err)
		}

		var empty IMSCC
		if reflect.DeepEqual(cc, empty) {
			t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
		}

		manifest, err := cc.ParseManifest()

		if err != nil {
			t.Errorf("Error parsing manifest: %v", err)
		}

		if reflect.TypeOf(manifest).Kind() != reflect.Struct {
			t.Errorf("Expecting struct type, got: %v", reflect.TypeOf(manifest).Kind())
		}
	}
}

func TestDump(t *testing.T) {

	cc := load(t, imscc_path)
	dump := cc.Dump()

	if reflect.TypeOf(dump).Kind() != reflect.Slice {
		t.Errorf("Expecting slice type, got: %v", reflect.TypeOf(dump).Kind())
	}

	if len(dump) == 0 {
		t.Error("Empty byte array returned!")
	}
}

func TestAsObject(t *testing.T) {

	cc := load(t, imscc_path)

	obj, err := cc.AsObject()

	if err != nil {
		t.Errorf("Error parsing the JSON: %v\n", err)
	}

	if reflect.TypeOf(obj).Kind() != reflect.Slice {
		t.Errorf("Expecting slice type, got: %v", reflect.TypeOf(obj).Kind())
	}

	if len(obj) == 0 {
		t.Error("Empty byte array returned!")
	}

}

func load(t *testing.T, p string) Cartridge {
	cc, err := Load(p)
	if err != nil {
		t.Errorf("could not load %s: %s", p, err)
	}
	return cc
}

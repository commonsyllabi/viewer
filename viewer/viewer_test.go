package viewer

import (
	"fmt"
	"reflect"
	"testing"
)

const imscc_path = "../test_files/canvas_large_1.3.imscc"

// have setup() to create tmp files, with e.g. different schema versions
// check go test tables

func TestLoadEmpty(t *testing.T) {
	var cc Cartridge //-- declare interface instance
	cc = NewIMSCC()  //-- assign underlying value
	_, err := cc.Load("")

	if err == nil {
		t.Fail()
	}
}

func TestParseManifest(t *testing.T) {
	var cc Cartridge
	cc = NewIMSCC()
	cc, _ = cc.Load(imscc_path)

	manifest, err := cc.ParseManifest()

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if reflect.TypeOf(manifest).Kind() != reflect.Struct {
		fmt.Printf("Expecting struct type, got: %v", reflect.TypeOf(manifest).Kind())
		t.Fail()
	}
}

func TestLoadCorrect(t *testing.T) {
	var cc Cartridge //-- declare interface instance
	cc = NewIMSCC()  //-- assign underlying value
	cc, err := cc.Load(imscc_path)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var empty IMSCC
	if reflect.DeepEqual(cc, empty) {
		fmt.Printf("Expecting struct type, got: %v", reflect.TypeOf(cc).Kind())
		t.Fail()
	}
}

func TestDump(t *testing.T) {
	var cc Cartridge //-- declare interface instance
	cc = NewIMSCC()
	cc, _ = cc.Load(imscc_path)

	dump := cc.Dump()

	if reflect.TypeOf(dump).Kind() != reflect.Slice {
		fmt.Printf("Expecting slice type, got: %v", reflect.TypeOf(dump).Kind())
		t.Fail()
	}

	if len(dump) == 0 {
		fmt.Println("Empty byte array returned!")
		t.Fail()
	}
}

func TestAsObject(t *testing.T) {
	var cc Cartridge
	cc = NewIMSCC()
	cc, _ = cc.Load(imscc_path)

	obj, err := cc.AsObject()

	if err != nil {
		fmt.Printf("Error parsing the JSON: %v\n", err)
		t.Fail()
	}

	if reflect.TypeOf(obj).Kind() != reflect.Slice {
		fmt.Printf("Expecting slice type, got: %v", reflect.TypeOf(obj).Kind())
		t.Fail()
	}

	if len(obj) == 0 {
		fmt.Println("Empty byte array returned!")
		t.Fail()
	}

}

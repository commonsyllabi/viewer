package viewer

import (
	"testing"
)

// have setup() to create tmp files, with e.g. different schema versions
// check go test tables
// `reflect.deepEqual` to check if two types are similar

func TestLoadFileEmpty(t *testing.T) {
	_, err := LoadFile("")

	if err == nil {
		t.Fail()
	}
}

func TestLoadFileCorrect(t *testing.T) {
	file, err := LoadFile("../test_files/canvas_large_1.3.imscc")

	if err != nil {
		t.Fail()
	}

	if file == "" {
		t.Fail()
	}
}

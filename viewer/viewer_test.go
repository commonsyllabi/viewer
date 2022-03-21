package viewer

import (
	"testing"
)

// have setup() to create tmp files, with e.g. different schema versions
// check go test tables
// `reflect.deepEqual` to check if two types are similar

func TestLoadEmpty(t *testing.T) {
	_, err := Load("")

	if err == nil {
		t.Fail()
	}
}

func TestLoadCorrect(t *testing.T) {
	file, err := Load("../test_files/canvas_large_1.3.imscc")

	if err != nil {
		t.Fail()
	}

	if file == "" {
		t.Fail()
	}
}

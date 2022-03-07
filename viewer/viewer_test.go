package viewer

import (
	"testing"
)

func TestLoadFileEmpty(t *testing.T) {
	err := LoadFile("")

	if err == nil {
		t.Fail()
	}
}

func TestLoadFileCorrect(t *testing.T) {
	err := LoadFile("../test_files/canvas_large_1.3.imscc")

	if err != nil {
		t.Fail()
	}
}

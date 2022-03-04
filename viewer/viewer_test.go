package viewer

import "testing"

func TestReadFile(t *testing.T) {
	err := ReadFile()

	if err != nil {
		t.Fatal(err)
	}
}

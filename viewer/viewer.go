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

func ParseManifest(manifest_path string) error {
	file, err := os.OpenFile(manifest_path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Printf("opened manifest %s\n", file.Name())

	bytesArray, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var manifest Manifest

	xml.Unmarshal(bytesArray, &manifest)

	s, _ := json.MarshalIndent(manifest, "", "\t")
	fmt.Printf("manifest struct: %s\n", s)

	fmt.Printf("title: %v\n", manifest.Metadata.LOM.General.Title.Value)
	return nil
}

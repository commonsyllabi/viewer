package api

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds port numbers, target directories
type Config struct {
	Port         string `yaml:"port"`
	TmpDir       string `yaml:"tmpDir"`
	UploadsDir   string `yaml:"uploadsDir"`
	FilesDir     string `yaml:"filesDir"`
	PublicDir    string `yaml:"publicDir"`
	TemplatesDir string `yaml:"templatesDir"`
	Debug        bool   `yaml:"debug"`
}

var conf Config

// defaults is called if there is an error opening and parsing the config file
func (c *Config) defaults() {
	c.Port = "2046"
	c.TmpDir = "/tmp/commonsyllabi"
	c.UploadsDir = "uploads"
	c.FilesDir = "files"
	c.PublicDir = "./www/public"
	c.TemplatesDir = "./internal/api/templates"
	c.Debug = true
}

// load tries to load a yaml file from disk, and marshals it
func (cc *Config) load(path string) error {
	cwd, _ := os.Getwd()
	content, err := ioutil.ReadFile(filepath.Join(cwd, path))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &cc)
	return err
}

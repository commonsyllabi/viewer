package api

import (
	"io/ioutil"
	"os"
	"path/filepath"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"gopkg.in/yaml.v2"
)

// Config holds port numbers, target directories
type Config struct {
	port       string `yaml:"port"`
	tmpDir     string `yaml:"tmpDir"`
	uploadsDir string `yaml:"uploadsDir"`
	filesDir   string `yaml:"filesDir"`
	publicDir  string `yaml:"publicDir"`
}

var conf Config

// defaults is called if there is an error opening and parsing the config file
func (c *Config) defaults() {
	c.port = "3046"
	c.tmpDir = "/tmp/commonsyllabi"
	c.uploadsDir = "uploads"
	c.filesDir = "files"
	c.publicDir = "./internal/www/public"
}

// load tries to load a yaml file from disk, and marshals it
func (cc *Config) load(path string) error {
	var c Config
	cwd, _ := os.Getwd()
	content, err := ioutil.ReadFile(filepath.Join(cwd, path))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &c)
	zero.Log.Debug().Msgf("%+v", c)
	return err
}

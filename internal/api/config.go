package api

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds port numbers, target directories
type Config struct {
	TmpDir       string `yaml:"tmpDir"`
	UploadsDir   string `yaml:"uploadsDir"`
	FilesDir     string `yaml:"filesDir"`
	PublicDir    string `yaml:"publicDir"`
	TemplatesDir string `yaml:"templatesDir"`
	FixturesDir  string `yaml:"FixturesDir"`
}

// DefaultConf is called if there is an error opening and parsing the config file
func (c *Config) DefaultConf() {
	c.TmpDir = "/tmp/cosyll/viewer"
	c.UploadsDir = "uploads"
	c.FilesDir = "files"
	c.PublicDir = os.Getenv("COSYLL_VIEWER_PUBLIC_DIR")
}

// LoadConf tries to load a yaml file from disk, and marshals it. Sensible defaults are provided, and loading a file overrides them
func (c *Config) LoadConf(path string) error {
	cwd, _ := os.Getwd()
	content, err := os.ReadFile(filepath.Join(cwd, path))
	if err != nil {
		c.DefaultConf()
		return err
	}

	err = yaml.Unmarshal(content, &c)
	if err != nil {
		c.DefaultConf()
		return err
	}
	return nil
}

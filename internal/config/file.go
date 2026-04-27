package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type URLs struct {
	Web string `yaml:"web"`
}

type File struct {
	URLs URLs `yaml:"urls"`
}

func (file *File) Validate() error {

	protocolAndPortOnly := func(raw string) bool {
		u, err := url.Parse(raw)
		if err != nil {
			return false
		}
		return u.Scheme != "" &&
			u.Hostname() != "" &&
			u.Port() == "" &&
			u.Path == "" &&
			u.RawQuery == "" &&
			u.Fragment == ""
	}

	if !protocolAndPortOnly(file.URLs.Web) {
		return fmt.Errorf("invalid URL: %s (expected format: https://example.com, http://localhost, etc.)", file.URLs.Web)
	}

	return nil
}

//go:embed file.yaml
var TemplateConfigFile []byte

func GetConfigFilePath() string {
	return filepath.Join(GetConfigDirPath(), "config.yaml")
}

func LoadFile() (*File, error) {
	return LoadFileFrom(GetConfigFilePath())
}

func LoadFileFrom(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found; have you run `projdocs init`?")
		}
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	defer f.Close()

	var cfg File
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &cfg, nil
}

func init() {
	var cfg File
	dec := yaml.NewDecoder(bytes.NewReader(TemplateConfigFile))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		panic(fmt.Sprintf("embedded config template does not match File struct: %v", err))
	}
}

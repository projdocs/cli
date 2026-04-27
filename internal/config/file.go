package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type URLs struct {
	Web      string `yaml:"web"`
	Supabase string `yaml:"supabase"`
}

type File struct {
	URLs URLs `yaml:"urls"`
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

package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/natefinch/atomic"
	"gopkg.in/yaml.v3"
)

type Presets map[string]Config

var PathConfig = filepath.Join(GetGoEnv("GOPATH"), ".go-env.yml")

func LoadFromConfig() Presets {

	file, err := os.OpenFile(PathConfig, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	var presets = make(Presets)
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(presets)

	if err != nil {
		if err == io.EOF {
			return presets
		}

		log.Fatalf("Failed to decode YAML: %v", err)
	}

	return presets
}

func (p *Presets) Keys() []string {
	keys := make([]string, 0, len(*p))
	for k := range *p {
		keys = append(keys, k)
	}

	return keys
}

func (p *Presets) Remove(name string) {
	delete(*p, name)
}

func (p *Presets) Persist() error {
	yamlData, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("can`t encode presets to yaml: %w", err)
	}

	reader := bytes.NewReader(yamlData)
	err = atomic.WriteFile(PathConfig, reader)
	if err != nil {
		return fmt.Errorf("can`t save config file: %w", err)
	}

	return nil
}

package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ManifestLoader struct {
	data []byte
}

var defaultManifests = []string{"cloudia.json", "cloudia.yaml", "cloudia.yml"}

///////////////////////////////////////////////////////////////////////////////////////////////
// Thanks to https://github.com/markpeek/packer/commit/5bf33a0e91b2318a40c42e9bf855dcc8dd4cdec5
///////////////////////////////////////////////////////////////////////////////////////////////

func (m *ManifestLoader) formatSyntaxError(syntaxError error) (err error) {
	syntax, ok := syntaxError.(*json.SyntaxError)
	if !ok {
		err = syntaxError
		return
	}

	data := m.data
	newline := []byte{'\x0a'}
	space := []byte{' '}

	start, end := bytes.LastIndex(data[:syntax.Offset], newline)+1, len(data)
	if idx := bytes.Index(data[start:], newline); idx >= 0 {
		end = start + idx
	}

	line, pos := bytes.Count(data[:start], newline)+1, int(syntax.Offset)-start-1
	err = fmt.Errorf("\nError in line %d: %s \n%s\n%s^",
		line, syntaxError, data[start:end], bytes.Repeat(space, pos))
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *ManifestLoader) unmarshalJSON() (Manifest, error) {
	var manifest Manifest
	err := json.Unmarshal(m.data, &manifest)
	if err != nil {
		err = m.formatSyntaxError(err)
		return nil, err
	}
	return manifest, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *ManifestLoader) unmarshalYAML() (Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal(m.data, &manifest)
	if err != nil {
		err = m.formatSyntaxError(err)
		return nil, err
	}
	return manifest, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func manifestFiles() []string {
	var result = []string(nil)
	if len(options.manifest) > 0 {
		result = []string{options.manifest}
	} else {
		result = defaultManifests
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *ManifestLoader) LoadFromFile(filename string) (Manifest, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	if ext == ".json" {
		return m.unmarshalJSON()
	} else if ext == ".yml" || ext == ".yaml" {
		return m.unmarshalYAML()
	} else if ext == "" {
		return m.unmarshalJSON()
	} else {
		return nil, errors.New("Unrecognized file extension")
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *ManifestLoader) LoadRaw(rawManifest string) (Manifest, error) {
	if len(rawManifest) > 0 {
		m.data = []byte(rawManifest)
		return m.unmarshalJSON()
	} else {
		for _, f := range manifestFiles() {
			if _, err := os.Stat(f); err == nil {
				return m.LoadFromFile(f)
			}
		}
	}

	return nil, StatusError{fmt.Errorf("no manifest found %v", manifestFiles()), 78}
}

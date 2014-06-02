package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/v1/yaml"
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
		return nil, StatusError{err, 65}
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
		return nil, StatusError{err, 65}
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
func (m *ManifestLoader) readFromFile(filename string) (Manifest, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, StatusError{err, 74}
	}

	ext := filepath.Ext(filename)
	if ext == ".json" {
		return unmarshalJSON(data)
	} else if ext == ".yml" || ext == ".yaml" {
		return unmarshalYAML(data)
	} else if ext == "" {
		return unmarshalJSON(data)
	} else {
		return nil, StatusError{errors.New("Unrecognized file extension"), 65}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *ManifestLoader) GetManifest(rawManifest string) (Manifest, error) {
	if len(rawManifest) > 0 {
		return m.unmarshalJSON([]byte(rawManifest))
	} else {
		for _, f := range manifestFiles() {
			if _, err := os.Stat(f); err == nil {
				return readFromFile(f)
			}
		}
	}

	return nil, StatusError{fmt.Errorf("no manifest found %v", manifestFiles()), 78}
}

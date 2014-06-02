package manifest

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
type Manifest struct {
	Containers Containers          `json:"containers" yaml:"containers"`
	Groups     map[string][]string `json:"groups" yaml:"groups"`
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *Manifest) GetTargetedContainers(group string) ([]string, error) {
	// If group is not given, all containers
	if len(group) == 0 {
		var names []string
		for i := 0; i < len(m.Containers); i++ {
			names = append(names, m.Containers[i].Name())
		}
		return names, nil
	}
	// Select specified group from listed groups
	for name, containers := range manifest.Groups {
		if name == group {
			return containers, nil
		}
	}
	// The group might just be a container reference itself
	for i := 0; i < len(manifest.Containers); i++ {
		if manifest.Containers[i].Name() == group {
			return append([]string{}, group), nil
		}
	}
	// Otherwise, fail verbosely
	return nil, StatusError{fmt.Errorf("no group nor container matching `%s`", group), 64}
}

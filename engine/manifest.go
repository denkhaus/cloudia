package engine

import (
	"fmt"
	"os"
)

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
type Manifest struct {
	Containers Containers          `json:"containers" yaml:"containers"`
	Groups     map[string][]string `json:"groups" yaml:"groups"`
	Nodes      map[string]string   `json:"nodes" yaml:"nodes"`
}

///////////////////////////////////////////////////////////////////////////////////////////////
// containerInGroup
///////////////////////////////////////////////////////////////////////////////////////////////
func containerInGroup(container Container, names []string) bool {
	for _, name := range names {
		if os.ExpandEnv(name) == container.Name() {
			return true
		}
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////
// GetContainersByNames
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *Manifest) GetContainersByNames(names []string) []Container {
	cont := m.Containers
	var filtered []Container
	for i := 0; i < cont.Count(); i++ {
		if containerInGroup(cont[i], names) {
			filtered = append(filtered, cont[i])
		}
	}
	return filtered
}

///////////////////////////////////////////////////////////////////////////////////////////////
//GetContainerNamesByGroup
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *Manifest) GetContainerNamesByGroup(group string) ([]string, error) {
	cont := m.Containers
	// If group is not given, all containers
	if len(group) == 0 {
		var names []string
		for i := 0; i < cont.Count(); i++ {
			names = append(names, cont[i].Name())
		}
		return names, nil
	}
	// Select specified group from listed groups
	for groupName, containerNames := range m.Groups {
		if groupName == group {
			return containerNames, nil
		}
	}
	// The group might just be a container reference itself
	for i := 0; i < len(m.Containers); i++ {
		if m.Containers[i].Name() == group {
			return append([]string{}, group), nil
		}
	}

	return nil, fmt.Errorf("no group nor container matching `%s`", group)
}

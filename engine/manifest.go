package engine

import (
	"os"
)

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
type Manifest struct {
	Containers []Container         `json:"containers" yaml:"containers"`
	Groups     map[string][]string `json:"groups" yaml:"groups"`
	Nodes      Nodes               `json:"nodes" yaml:"nodes"`
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
	var filtered []Container
	for _, cnt := range m.Containers {
		if containerInGroup(cnt, names) {
			filtered = append(filtered, cnt)
		}
	}

	return filtered
}

///////////////////////////////////////////////////////////////////////////////////////////////
//GetContainerNamesByGroup
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *Manifest) GetContainerNamesByGroup(group string) []string {
	// If group is not given, all containers
	if len(group) == 0 {
		var names []string
		for _, cnt := range m.Containers {
			names = append(names, cnt.Name())
		}
		return names
	}
	// Select specified group from listed groups
	for groupName, containerNames := range m.Groups {
		if groupName == group {
			return containerNames
		}
	}
	// The group might just be a container reference itself
	var names []string
	for _, cnt := range m.Containers {
		if cnt.Name() == group {
			names = append(names, group)
		}
	}

	return names
}

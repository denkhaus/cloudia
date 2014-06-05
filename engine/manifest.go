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
	cont.ForAll(func(cnt Container) error {
		if containerInGroup(cnt, names) {
			filtered = append(filtered, cnt)
		}
		return nil
	})

	return filtered
}

///////////////////////////////////////////////////////////////////////////////////////////////
//GetContainerNamesByGroup
///////////////////////////////////////////////////////////////////////////////////////////////
func (m *Manifest) GetContainerNamesByGroup(group string) []string {
	cont := m.Containers
	// If group is not given, all containers
	if len(group) == 0 {
		var names []string
		cont.ForAll(func(cnt Container) error {
			names = append(names, cnt.Name())
			return nil
		})
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
	cont.ForAll(func(cnt Container) error {
		if cnt.Name() == group {
			names = append(names, group)
		}
		return nil
	})

	return names
}

package engine

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
	"strconv"
	"strings"
)

type RunParameters struct {
	RawCidfile     string      `json:"cidfile" yaml:"cidfile"`
	CpuShares      int         `json:"cpu-shares" yaml:"cpu-shares"`
	RawDns         []string    `json:"dns" yaml:"dns"`
	RawEntrypoint  []string    `json:"entrypoint" yaml:"entrypoint"`
	RawEnv         []string    `json:"env" yaml:"env"`
	RawEnvFile     string      `json:"env-file" yaml:"env-file"`
	RawExpose      []string    `json:"expose" yaml:"expose"`
	RawHostname    string      `json:"hostname" yaml:"hostname"`
	Interactive    bool        `json:"interactive" yaml:"interactive"`
	RawLink        []string    `json:"link" yaml:"link"`
	RawLxcConf     []string    `json:"lxc-conf" yaml:"lxc-conf"`
	RawMemory      string      `json:"memory" yaml:"memory"`
	RawMemorySwap  string      `json:"memory-swap" yaml:"memory-swap"`
	RawNet         string      `json:"net" yaml:"net"`
	Privileged     bool        `json:"privileged" yaml:"privileged"`
	RawPorts       []string    `json:"ports" yaml:"ports"`
	PublishAll     bool        `json:"publish-all" yaml:"publish-all"`
	Rm             bool        `json:"rm" yaml:"rm"`
	Tty            bool        `json:"tty" yaml:"tty"`
	RawUser        string      `json:"user" yaml:"user"`
	RawBinds       []string    `json:"binds" yaml:"binds"`
	RawVolumesFrom string      `json:"volumes-from" yaml:"volumes-from"`
	RawWorkdir     string      `json:"workdir" yaml:"workdir"`
	RawCmd         interface{} `json:"cmd" yaml:"cmd"`
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Cidfile() string {
	return os.ExpandEnv(r.RawCidfile)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Dns() []string {
	var dns []string
	for _, rawDns := range r.RawDns {
		dns = append(dns, os.ExpandEnv(rawDns))
	}
	return dns
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Entrypoint() []string {
	var ep []string
	for _, rawEp := range r.RawEntrypoint {
		ep = append(ep, os.ExpandEnv(rawEp))
	}
	return ep
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Env() []string {
	var env []string
	for _, rawEnv := range r.RawEnv {
		env = append(env, os.ExpandEnv(rawEnv))
	}
	return env
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) EnvFile() string {
	return os.ExpandEnv(r.RawEnvFile)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Expose() []string {
	var expose []string
	for _, rawExpose := range r.RawExpose {
		expose = append(expose, os.ExpandEnv(rawExpose))
	}
	return expose
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Hostname() string {
	return os.ExpandEnv(r.RawHostname)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Link() []string {
	var link []string
	for _, rawLink := range r.RawLink {
		link = append(link, os.ExpandEnv(rawLink))
	}
	return link
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Memory() string {
	return os.ExpandEnv(r.RawMemory)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) MemorySwap() string {
	return os.ExpandEnv(r.RawMemorySwap)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Net() string {
	// Default to bridge
	if len(r.RawNet) == 0 {
		return "bridge"
	} else {
		return os.ExpandEnv(r.RawNet)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) PortMappings() []string {
	var ports []string
	for _, rawPort := range r.RawPorts {
		var expPorts []string
		for _, port := range strings.Split(rawPort, ":") {
			expPorts = append(expPorts, os.ExpandEnv(port))
		}
		ports = append(ports, strings.Join(expPorts, ":"))
	}
	return ports
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) PortBindings() map[docker.Port][]docker.PortBinding {
	bindings := make(map[docker.Port][]docker.PortBinding)

	for _, mapping := range r.PortMappings() {
		var bndngs []docker.PortBinding
		tok := strings.Split(mapping, ":")
		bndngs = append(bndngs, docker.PortBinding{HostPort: tok[0]})
		bindings[docker.Port(tok[1]+"/tcp")] = bndngs
	}

	return bindings
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) ExposedPorts() map[docker.Port]struct{} {
	type param struct{}
	ports := make(map[docker.Port]struct{})

	for port, _ := range r.PortBindings() {
		ports[port] = param{}
	}

	return ports
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) User() string {
	return os.ExpandEnv(r.RawUser)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Binds() []string {
	var binds []string
	for _, rawBind := range r.RawBinds {
		var expPaths []string
		for _, path := range strings.Split(rawBind, ":") {
			expPaths = append(expPaths, os.ExpandEnv(path))
		}
		binds = append(binds, strings.Join(expPaths, ":"))
	}
	return binds
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Volumes() map[string]struct{} {
	type param struct{}
	var volumes = make(map[string]struct{})
	for _, bind := range r.Binds() {
		tok := strings.Split(bind, ":")
		volumes[tok[1]] = param{}
	}
	return volumes
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) VolumesFrom() string {
	return os.ExpandEnv(r.RawVolumesFrom)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) WorkingDir() string {
	return os.ExpandEnv(r.RawWorkdir)
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) LxcConf() []docker.KeyValuePair {
	var cnf []docker.KeyValuePair
	for _, rawConf := range r.RawLxcConf {
		tok := strings.Split(rawConf, ":")
		kv := docker.KeyValuePair{
			Key:   os.ExpandEnv(tok[0]),
			Value: os.ExpandEnv(tok[1]),
		}

		cnf = append(cnf, kv)
	}

	return cnf
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) Cmd() ([]string, error) {
	var cmd []string
	if r.RawCmd != nil {
		switch rawCmd := r.RawCmd.(type) {
		case string:
			if len(rawCmd) > 0 {
				return append(cmd, os.ExpandEnv(rawCmd)), nil
			}
		case []interface{}:
			cmds := make([]string, len(rawCmd))
			for i, v := range rawCmd {
				cmds[i] = os.ExpandEnv(v.(string))
			}
			return append(cmd, cmds...), nil
		}
	}

	return nil, errors.New("cmd is of unknown type!")
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) CreateDockerHostConfig() (*docker.HostConfig, error) {
	config := &docker.HostConfig{}

	config.Privileged = r.Privileged
	config.PublishAllPorts = r.PublishAll
	config.ContainerIDFile = r.Cidfile()
	config.PortBindings = r.PortBindings()
	config.LxcConf = r.LxcConf()

	config.Dns = append(config.Dns, r.Dns()...)
	config.Links = append(config.Links, r.Link()...)
	config.Binds = append(config.Binds, r.Binds()...)

	return config, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) CreateDockerConfig() (*docker.Config, error) {
	config := &docker.Config{}

	//TODO
	//Domainname      string
	//PortSpecs       []string
	//OpenStdin       bool
	//StdinOnce       bool
	//Image           string
	//NetworkDisabled bool

	//TODO defaults
	config.AttachStdout = true
	config.AttachStdin = true
	config.AttachStderr = false

	config.Tty = r.Tty
	config.User = r.User()
	config.Hostname = r.Hostname()
	config.VolumesFrom = r.VolumesFrom()
	config.WorkingDir = r.WorkingDir()
	config.ExposedPorts = r.ExposedPorts()
	config.CpuShares = int64(r.CpuShares)
	config.Volumes = r.Volumes()

	config.Dns = append(config.Dns, r.Dns()...)
	config.Env = append(config.Env, r.Env()...)
	config.Entrypoint = append(config.Env, r.Entrypoint()...)

	if cmds, err := r.Cmd(); err != nil {
		config.Cmd = append(config.Cmd, cmds...)
	} else {
		return nil, fmt.Errorf("Run parameters error:: Errror while parsing cmd:: %v", err)
	}

	// Memory
	if len(r.Memory()) > 0 {
		if mem, err := strconv.ParseInt(r.Memory(), 10, 64); err == nil {
			config.Memory = mem
		} else {
			return nil, fmt.Errorf("Run parameters error:: Unable to convert memory param to int64 %v", err)
		}

	}

	// MemorySwap
	if len(r.MemorySwap()) > 0 {
		if swap, err := strconv.ParseInt(r.MemorySwap(), 10, 64); err == nil {
			config.MemorySwap = swap
		} else {
			return nil, fmt.Errorf("Run parameters error:: Unable to convert memory-swap param to int64 %v", err)
		}

	}

	return config, nil
}

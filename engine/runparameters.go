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
	Detach         bool        `json:"detach" yaml:"detach"`
	RawDns         []string    `json:"dns" yaml:"dns"`
	RawEntrypoint  string      `json:"entrypoint" yaml:"entrypoint"`
	RawEnv         []string    `json:"env" yaml:"env"`
	RawEnvFile     string      `json:"env-file" yaml:"env-file"`
	RawExpose      []string    `json:"expose" yaml:"expose"`
	RawHostname    string      `json:"hostname" yaml:"hostname"`
	Interactive    bool        `json:"interactive" yaml:"interactive"`
	RawLink        []string    `json:"link" yaml:"link"`
	RawLxcConf     []string    `json:"lxc-conf" yaml:"lxc-conf"`
	RawMemory      string      `json:"memory" yaml:"memory"`
	RawNet         string      `json:"net" yaml:"net"`
	Privileged     bool        `json:"privileged" yaml:"privileged"`
	RawPublish     []string    `json:"publish" yaml:"publish"`
	PublishAll     bool        `json:"publish-all" yaml:"publish-all"`
	Rm             bool        `json:"rm" yaml:"rm"`
	Tty            bool        `json:"tty" yaml:"tty"`
	RawUser        string      `json:"user" yaml:"user"`
	RawVolumes     []string    `json:"volume" yaml:"volume"`
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
func (r *RunParameters) Entrypoint() string {
	return os.ExpandEnv(r.RawEntrypoint)
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
func (r *RunParameters) LxcConf() []string {
	var lxcConf []string
	for _, rawLxcConf := range r.RawLxcConf {
		lxcConf = append(lxcConf, os.ExpandEnv(rawLxcConf))
	}
	return lxcConf
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
func (r *RunParameters) Publish() []string {
	var publish []string
	for _, rawPublish := range r.RawPublish {
		publish = append(publish, os.ExpandEnv(rawPublish))
	}
	return publish
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
func (r *RunParameters) Volumes() []string {
	var volumes []string
	for _, rawVolume := range r.RawVolumes {
		paths := strings.Split(rawVolume, ":")
		volumes = append(volumes, os.ExpandEnv(strings.Join(paths, ":")))
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
func (r *RunParameters) CreateDockerConfig() (*docker.HostConfig, error) {
	config := &docker.HostConfig{}
	//		if len(container.Run.Cidfile()) > 0 {
	//			args = append(args, "--cidfile", container.Run.Cidfile())
	//		}

	//		// Detach
	//		if container.Run.Detach {
	//			args = append(args, "--detach")
	//		}

	//		// Entrypoint
	//		if len(container.Run.Entrypoint()) > 0 {
	//			args = append(args, "--entrypoint", container.Run.Entrypoint())
	//		}

	//		// Env file
	//		if len(container.Run.EnvFile()) > 0 {
	//			args = append(args, "--env-file", container.Run.EnvFile())
	//		}
	//		// Expose
	//		for _, expose := range container.Run.Expose() {
	//			args = append(args, "--expose", expose)
	//		}

	//		// Interactive
	//		if container.Run.Interactive {
	//			args = append(args, "--interactive")
	//		}
	//		// Link
	//		for _, link := range container.Run.Link() {
	//			args = append(args, "--link", link)
	//		}
	//		// LxcConf
	//		for _, lxcConf := range container.Run.LxcConf() {
	//			args = append(args, "--lxc-conf", lxcConf)
	//		}

	//		// Net
	//		if container.Run.Net() != "bridge" {
	//			args = append(args, "--net", container.Run.Net())
	//		}

	//		// Privileged
	//		if container.Run.Privileged {
	//			args = append(args, "--privileged")
	//		}
	//		// Publish
	//		for _, port := range container.Run.Publish() {
	//			args = append(args, "--publish", port)
	//		}
	//		// PublishAll
	//		if container.Run.PublishAll {
	//			args = append(args, "--publish-all")
	//		}
	//		// Rm
	//		if container.Run.Rm {
	//			args = append(args, "--rm")
	//		}
	//		// Tty
	//		if container.Run.Tty {
	//			args = append(args, "--tty")
	//		}

	return config
}

///////////////////////////////////////////////////////////////////////////////////////////////
//
///////////////////////////////////////////////////////////////////////////////////////////////
func (r *RunParameters) CreateDockerConfig() (*docker.Config, error) {

	config := &docker.Config{}

	// CPU shares
	if r.CpuShares > 0 {
		config.CpuShares = int64(r.CpuShares)
	}
	// Dns
	for _, dns := range r.Dns() {
		config.Dns = append(config.Dns, dns)
	}
	// Env
	for _, env := range r.Env() {
		config.Env = append(config.Env, env)
	}
	// Host
	if len(r.Hostname()) > 0 {
		config.Hostname = r.Hostname()
	}
	// Memory
	if len(r.Memory()) > 0 {
		if mem, err := strconv.ParseInt(r.Memory(), 10, 64); err == nil {
			config.Memory = mem
		} else {
			return nil, fmt.Errorf("Run parameters error:: Unable to convert memory param to int64 %v", err)
		}

	}
	// User
	if len(r.User()) > 0 {
		config.User = r.User()
	}

	// TODO Volumes
	//for _, volume := range cnt.Run.Volumes() {
	//	config.Volumes = append(config.Volumes, volume)
	//}

	// VolumesFrom
	if len(r.VolumesFrom()) > 0 {
		config.VolumesFrom = r.VolumesFrom()
	}
	// WorkingDir
	if len(r.WorkingDir()) > 0 {
		config.WorkingDir = r.WorkingDir()
	}

	// Command
	if cmds, err := r.Cmd(); err != nil {
		for _, cmd := range cmds {
			config.Cmd = append(config.Cmd, cmd)
		}
	} else {
		return nil, fmt.Errorf("Run parameters error:: Errror while parsing cmd:: %v", err)
	}

	return config, nil
}

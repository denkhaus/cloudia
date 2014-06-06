package engine

import (
	"errors"
	"os"
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
	for _, rawVolume := range r.RawVolume {
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

	//TODO more verbose error notification
	return nil, errors.New("cmd is of unknown type!")
}

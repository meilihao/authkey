// from https://github.com/pingcap/tidb/blob/master/util/versioninfo/versioninfo.go
package util

import (
	"encoding/json"
	"runtime"
	"strings"
)

var (
	version        = ""
	gitBranch      = ""
	gitHash        = ""
	buildTimestamp = ""

	gitStateDirtySuffix = "-dirty"
)

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitBranch is the brach of the git tree.
	GitBranch string `json:"git_brach,omitempty"`
	// GitHash is the git sha1.
	GitHash string `json:"git_hash,omitempty"`
	// GitState is the state of the git tree
	GitState string `json:"git_state,omitempty"`
	// BuildTimestate build time
	BuildTimestamp string `json:"build_timestamp,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

// Get returns build info
func GetBuildInfo() BuildInfo {
	i := BuildInfo{
		Version:        version,
		GitBranch:      gitBranch,
		GitHash:        gitHash,
		GitState:       "clean",
		BuildTimestamp: buildTimestamp,
		GoVersion:      runtime.Version(),
	}

	if strings.HasSuffix(i.Version, gitStateDirtySuffix) {
		i.Version = strings.TrimSuffix(i.Version, gitStateDirtySuffix)
		i.GitState = "dirty"
	}

	return i
}

func (i BuildInfo) String() string {
	data, _ := json.Marshal(i)

	return string(data)
}

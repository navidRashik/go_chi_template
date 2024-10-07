package version

import (
	"fmt"
	"runtime/debug"

	"example_project/assets"
)

func GetBuildVersion() string {
	var revision string
	var modified bool
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}
	if revision == "" {
		return "unavailable"
	}
	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}
	return revision
}

func GetVersion() string {
	dat, err := assets.EmbeddedFiles.ReadFile("VERSION")
	if err != nil {
		return ""
	}
	return string(dat)
}

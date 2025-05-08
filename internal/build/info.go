package build

import (
	"fmt"
	"runtime"
)

// Эти переменные будут перезаписаны при сборке через -ldflags.
var (
	version   string
	commit    string
	buildDate string
)

type Info struct {
	Version   string
	Commit    string
	BuildDate string
	GoVersion string
	OSArch    string
}

// GetWithBuildFlags возвращает BuildInfo с данными из -ldflags.
func GetWithBuildFlags() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		OSArch:    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

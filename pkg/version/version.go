package version

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Version provided by compile time -ldflags.
var (
	defaultContent = "was not build properly"
	VersionTag     = defaultContent
	Branch         = defaultContent
	BuildDate      = defaultContent
)

// Version holds information about the build
type Version struct {
	Version, Branch, GoVersion, Platform string
	BuildDate                            time.Time
}

// Get returns the version information
func Get() Version {
	v := Version{
		Version:   VersionTag,
		Branch:    Branch,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	i, err := strconv.ParseInt(BuildDate, 10, 64)
	if err != nil {
		panic(fmt.Errorf("Error parsing build time: %v", err))
	}
	v.BuildDate = time.Unix(i, 0).UTC()
	return v
}

// PrintText prints the version informations to stdout
func PrintText() {
	v := Get()
	fmt.Printf("Version:\t%s\n", v.Version)
	fmt.Printf("BuildDate:\t%s\n", v.BuildDate.Local().Format(time.UnixDate))
	fmt.Printf("Branch:\t\t%s\n", v.Branch)
	fmt.Printf("Go Version:\t%s\n", v.GoVersion)
	fmt.Printf("Platform:\t%s\n", v.Platform)
	os.Exit(0)
}

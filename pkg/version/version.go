package version

import "fmt"


var (
	// Commit holds the git commit hash
	Commit string
	// BuildDate holds the date of the build
	BuildDate string


)
func Version() string {
    return fmt.Sprintf("%s_%s", Commit, BuildDate) 
}

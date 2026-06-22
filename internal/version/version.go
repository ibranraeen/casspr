package version

import (
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

// Build-time parameters set via -ldflags.

var (
	Version = "devel"
	Commit  = "unknown"
	// BuildID is a unique identifier for this build. For release builds it
	// equals Commit; for development builds (go run / go build without
	// ldflags) it is derived from the executable's modification time, which
	// changes on every recompilation.
	BuildID = ""
)

var (
	// Matches Go pseudo-version suffixes like -0.20260619214928-1849ddd1abcd
	pseudoVersionRegex = regexp.MustCompile(`-(?:[a-zA-Z0-9-.]*\.)?[0-9]{14}-[0-9a-fA-F]+$`)
	// Matches git describe suffixes like -2-g1849ddd1
	gitDescribeRegex = regexp.MustCompile(`-[0-9]+-g[0-9a-fA-F]+$`)
	// Matches raw commit hash suffix like -1849ddd1 (minimum 8 hex characters)
	commitSuffixRegex = regexp.MustCompile(`-[0-9a-fA-F]{8,40}$`)
)

func prunePseudoVersion(v string) string {
	// Strip build metadata/dirty suffix (e.g. +dirty)
	v = strings.Split(v, "+")[0]
	v = pseudoVersionRegex.ReplaceAllString(v, "")
	v = gitDescribeRegex.ReplaceAllString(v, "")
	v = commitSuffixRegex.ReplaceAllString(v, "")
	return v
}

// A user may install crush using `go install github.com/ibranraeen/casspr@latest`.
// without -ldflags, in which case the version above is unset. As a workaround
// we use the embedded build version that *is* set when using `go install` (and
// is only set for `go install` and not for `go build`).
func init() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		mainVersion := info.Main.Version
		if mainVersion != "" && mainVersion != "(devel)" {
			Version = mainVersion
		}
	}

	Version = prunePseudoVersion(Version)

	// Derive BuildID when not set via ldflags.
	if BuildID == "" {
		BuildID = deriveBuildID()
	}
}

// deriveBuildID uses the running executable's modification time as a unique
// build fingerprint. This changes on every recompilation (including `go run`),
// making it reliable for detecting stale servers during development.
func deriveBuildID() string {
	exe, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	fi, err := os.Stat(exe)
	if err != nil {
		return "unknown"
	}
	return strconv.FormatInt(fi.ModTime().UnixNano(), 36)
}

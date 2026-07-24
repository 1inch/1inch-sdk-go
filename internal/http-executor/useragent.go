package http_executor

import (
	"regexp"
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/1inch/1inch-sdk-go/v4"

// semverRegex is the official semver.org grammar with the leading "v" Go module
// versions carry. Every version the Go toolchain resolves (release tags,
// prerelease tags, pseudo-versions) matches it; toolchain placeholders like
// "(devel)" and empty strings from local path replacements do not.
var semverRegex = regexp.MustCompile(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

var userAgent = buildUserAgent()

func buildUserAgent() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		info = nil
	}
	return "1inch-dev-portal-client-go:" + sdkVersionFromBuildInfo(info)
}

// sdkVersionFromBuildInfo resolves the SDK version the running binary was built
// against and reports it verbatim when it is valid semver, so the header stays
// accurate across releases without manual updates and never drops version data.
// Versions that are not versions at all (a "(devel)" placeholder, or the empty
// version of a local path replacement) report as "unknown".
func sdkVersionFromBuildInfo(info *debug.BuildInfo) string {
	if info == nil {
		return "unknown"
	}

	for _, dep := range info.Deps {
		if dep.Path != modulePath {
			continue
		}
		// A replace directive swaps in different code, so the replacement's
		// version describes what actually runs, whatever module path it has
		effective := dep
		if dep.Replace != nil {
			effective = dep.Replace
		}
		if semverRegex.MatchString(effective.Version) {
			return effective.Version
		}
		return "unknown"
	}

	// When built from within this repository the module is Main, not a dependency
	if strings.HasPrefix(info.Main.Path, modulePath) && semverRegex.MatchString(info.Main.Version) {
		return info.Main.Version
	}

	return "unknown"
}

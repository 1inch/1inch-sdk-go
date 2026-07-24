package http_executor

import (
	"regexp"
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/1inch/1inch-sdk-go/v4"

// releaseVersionRegex matches exact release tags (vX.Y.Z). Pseudo-versions
// (v4.1.1-0.20260801120000-abcdef123456) and toolchain placeholders ((devel))
// are excluded so the header only ever reports a published release.
var releaseVersionRegex = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

var userAgent = buildUserAgent()

// buildUserAgent derives the SDK version from the consumer's build info so the
// User-Agent header stays accurate across releases without manual updates. Only
// exact release versions are reported; builds from unreleased commits report
// "unknown"
func buildUserAgent() string {
	version := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == modulePath && releaseVersionRegex.MatchString(dep.Version) {
				version = dep.Version
				break
			}
		}
		// When built from within this repository the module is Main, not a dependency
		if version == "unknown" && strings.HasPrefix(info.Main.Path, modulePath) && releaseVersionRegex.MatchString(info.Main.Version) {
			version = info.Main.Version
		}
	}
	return "1inch-dev-portal-client-go:" + version
}

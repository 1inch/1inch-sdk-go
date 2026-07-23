package http_executor

import (
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/1inch/1inch-sdk-go/v4"

var userAgent = buildUserAgent()

// buildUserAgent derives the SDK version from the consumer's build info so the
// User-Agent header stays accurate across releases without manual updates
func buildUserAgent() string {
	version := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == modulePath {
				version = dep.Version
				break
			}
		}
		// When built from within this repository the module is Main, not a dependency
		if version == "unknown" && strings.HasPrefix(info.Main.Path, modulePath) && info.Main.Version != "" {
			version = info.Main.Version
		}
	}
	return "1inch-dev-portal-client-go:" + version
}

package http_executor

import (
	"github.com/1inch/1inch-sdk-go/v5/internal/version"
)

// userAgent identifies the SDK and its release version on every API request.
// The version comes from the release-managed constant, so the header always
// names the release the code shipped in.
var userAgent = "1inch-dev-portal-client-go:" + version.Version

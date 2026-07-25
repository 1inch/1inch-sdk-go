// Package version holds the SDK's single source-of-truth version.
package version

// Version is the most recent published release of this SDK and is reported in
// the User-Agent header of every API request. The "Release new version"
// workflow owns this file: it writes the new version and commits before
// tagging, so the value always matches the release tag. Never edit by hand.
// Exact release form only: vMAJOR.MINOR.PATCH.
const Version = "v4.0.1"

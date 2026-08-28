// Package version holds the SDK's single source-of-truth version.
package version

// DO NOT HAND-EDIT. The "Release new version" workflow owns this value: it
// writes the new version and commits before tagging, so the value always
// matches the release tag, and PR CI rejects manual changes.
//
// Version is the most recent published release of this SDK and is reported in
// the User-Agent header of every API request. Exact release form only:
// vMAJOR.MINOR.PATCH.
const Version = "v4.2.0"

package pireflect

import (
	"net/url"
	"regexp"
)

// Match from beginning of the line until the first `vendor/` (non-greedy)
var vendorRe = regexp.MustCompile("^.*?/vendor/")

// sanitizeFuncName makes the function name suitable for logging display. It removes
// url-encoded elements from the `dot.git` package names and shortens the
// vendored paths.
func sanitizeFuncName(function string) string {
	// Use the stdlib to un-escape any package import paths which can happen
	// in the case of the "dot-git" postfix. Seems like a bug in stdlib =/
	if unescaped, err := url.QueryUnescape(function); err == nil {
		function = unescaped
	}

	// strip everything prior to the vendor
	return vendorRe.ReplaceAllString(function, "vendor/")
}

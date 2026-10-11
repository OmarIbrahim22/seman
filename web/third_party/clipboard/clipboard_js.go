//go:build js

package clipboard

import "errors"

// The browser demo has no synchronous clipboard access.
var errUnsupported = errors.New("clipboard: not available in the browser build")

func readAll() (string, error) { return "", errUnsupported }

func writeAll(string) error { return errUnsupported }

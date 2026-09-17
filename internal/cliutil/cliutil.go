// Package cliutil provides shared helper functions used across
// multiple CLI subcommands in cmd/gaze and the report pipeline
// in internal/aireport. These helpers eliminate duplication of
// common patterns such as format validation and JSON capture.
package cliutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ValidateFormat checks that format is one of the supported output
// formats. The base set is always "text" and "json". Additional
// formats may be supplied to extend the allowlist for commands that
// support them (e.g. "html" for analyze). Callers that pass no extra
// values retain the exact text/json acceptance and error behavior.
func ValidateFormat(format string, extra ...string) error {
	if format == "text" || format == "json" {
		return nil
	}
	for _, e := range extra {
		if format == e {
			return nil
		}
	}
	if len(extra) == 0 {
		return fmt.Errorf("invalid format %q: must be 'text' or 'json'", format)
	}
	// Build a human-readable list preserving argument order.
	allowed := "'text', 'json'"
	for i, e := range extra {
		if i == len(extra)-1 {
			allowed += ", or '" + e + "'"
		} else {
			allowed += ", '" + e + "'"
		}
	}
	return fmt.Errorf("invalid format %q: must be %s", format, allowed)
}

// CaptureJSON calls fn with a buffer as the writer, then returns the
// buffer contents as a json.RawMessage. If fn returns an error, it is
// propagated and the message is nil.
func CaptureJSON(fn func(w io.Writer) error) (json.RawMessage, error) {
	var buf bytes.Buffer
	if err := fn(&buf); err != nil {
		return nil, err
	}
	return json.RawMessage(buf.Bytes()), nil
}

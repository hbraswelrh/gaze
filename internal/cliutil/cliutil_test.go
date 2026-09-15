package cliutil

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		extra   []string // optional additional allowed formats
		wantErr bool
		errSub  string // substring expected in error message
	}{
		// Base set (no extra) - unchanged behavior.
		{name: "valid text", format: "text", wantErr: false},
		{name: "valid json", format: "json", wantErr: false},
		{name: "invalid csv", format: "csv", wantErr: true, errSub: `"csv"`},
		{name: "empty string", format: "", wantErr: true, errSub: `""`},
		{name: "error message format", format: "xml", wantErr: true, errSub: "must be 'text' or 'json'"},

		// html rejected without extra - proves D4 scoping.
		{name: "html rejected without extra", format: "html", wantErr: true, errSub: `"html"`},

		// html accepted when explicitly supplied as extra.
		{name: "html accepted with extra", format: "html", extra: []string{"html"}, wantErr: false},

		// Base formats still accepted when extra is supplied.
		{name: "text still valid with extra", format: "text", extra: []string{"html"}, wantErr: false},
		{name: "json still valid with extra", format: "json", extra: []string{"html"}, wantErr: false},

		// Unsupported format still rejected even with extra.
		{name: "csv rejected with extra html", format: "csv", extra: []string{"html"}, wantErr: true, errSub: `"csv"`},

		// Error message lists all allowed formats when extra is supplied.
		{name: "error with extra lists all formats", format: "csv", extra: []string{"html"}, wantErr: true, errSub: "must be 'text', 'json', or 'html'"},

		// Legacy error message preserved when no extra is supplied.
		{name: "error without extra is legacy text", format: "csv", wantErr: true, errSub: "must be 'text' or 'json'"},

		// Multiple extra values.
		{name: "second extra accepted", format: "yaml", extra: []string{"html", "yaml"}, wantErr: false},

		// Multiple extra values in error message.
		{name: "error with two extras lists all", format: "bin", extra: []string{"html", "yaml"}, wantErr: true, errSub: "must be 'text', 'json', 'html', or 'yaml'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.format, tt.extra...)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidateFormat(%q, %v) = nil, want error", tt.format, tt.extra)
				}
				if tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
					t.Errorf("error %q does not contain %q", err, tt.errSub)
				}
			} else {
				if err != nil {
					t.Fatalf("ValidateFormat(%q, %v) = %v, want nil", tt.format, tt.extra, err)
				}
			}
		})
	}
}

var errSentinel = errors.New("sentinel error")

func TestCaptureJSON(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(w io.Writer) error
		wantNil bool
		wantStr string
		wantErr bool
	}{
		{
			name: "success path",
			fn: func(w io.Writer) error {
				_, err := fmt.Fprint(w, `{"key":"value"}`)
				return err
			},
			wantStr: `{"key":"value"}`,
		},
		{
			name: "error propagation preserves identity",
			fn: func(_ io.Writer) error {
				return errSentinel
			},
			wantNil: true,
			wantErr: true,
		},
		{
			name: "empty output returns nil message",
			fn: func(_ io.Writer) error {
				return nil
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CaptureJSON(tt.fn)
			if tt.wantErr {
				if err == nil {
					t.Fatal("CaptureJSON() = nil error, want error")
				}
				if !errors.Is(err, errSentinel) {
					t.Errorf("CaptureJSON() error = %v, want %v (identity preserved)", err, errSentinel)
				}
				if got != nil {
					t.Fatalf("CaptureJSON() returned non-nil message on error: %s", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("CaptureJSON() error = %v, want nil", err)
			}
			if tt.wantNil && got != nil {
				t.Fatalf("CaptureJSON() = %s, want nil", got)
			}
			if !tt.wantNil {
				if got == nil {
					t.Fatal("CaptureJSON() = nil, want non-nil json.RawMessage")
				}
				if string(got) != tt.wantStr {
					t.Errorf("CaptureJSON() = %q, want %q", string(got), tt.wantStr)
				}
				if tt.wantStr == "" && len(got) != 0 {
					t.Errorf("CaptureJSON() len = %d, want 0 for empty output", len(got))
				}
			}
		})
	}
}

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeInputLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		value string
		valid bool
	}{
		{name: "empty", input: "", value: "", valid: false},
		{name: "whitespace", input: "   \t", value: "", valid: false},
		{name: "comment", input: "  # skip", value: "", valid: false},
		{name: "comment_with_crlf", input: "\r\n# skip\r\n", value: "", valid: false},
		{name: "nul_byte_string", input: "\x00", value: "\x00", valid: true},
		{name: "value", input: "  example.org  ", value: "example.org", valid: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, valid := normalizeInputLine(tc.input)
			if value != tc.value || valid != tc.valid {
				t.Fatalf("normalizeInputLine(%q) = (%q, %t), want (%q, %t)", tc.input, value, valid, tc.value, tc.valid)
			}
		})
	}
}

func TestRequireReadableFileFromEnv(t *testing.T) {
	const envName = "IP_ENRICH_TEST_PATH"
	t.Cleanup(func() {
		_ = os.Unsetenv(envName)
	})

	t.Run("unset_env", func(t *testing.T) {
		_ = os.Unsetenv(envName)

		_, err := requireReadableFileFromEnv(envName)
		if err == nil {
			t.Fatal("expected error when env var is not set")
		}

		if !strings.Contains(err.Error(), "must be set") {
			t.Fatalf("unexpected error when env var is not set: %q", err.Error())
		}
	})

	t.Run("missing_file", func(t *testing.T) {
		missingPath := filepath.Join(t.TempDir(), "does-not-exist.mmdb")
		if err := os.Setenv(envName, missingPath); err != nil {
			t.Fatalf("failed to set env var: %v", err)
		}

		_, err := requireReadableFileFromEnv(envName)
		if err == nil {
			t.Fatal("expected error for missing file path")
		}

		if !strings.Contains(err.Error(), "cannot access") {
			t.Fatalf("unexpected error for missing file path: %q", err.Error())
		}
	})

	t.Run("directory_instead_of_file", func(t *testing.T) {
		dirPath := t.TempDir()
		if err := os.Setenv(envName, dirPath); err != nil {
			t.Fatalf("failed to set env var: %v", err)
		}

		_, err := requireReadableFileFromEnv(envName)
		if err == nil {
			t.Fatal("expected error for directory path")
		}

		if !strings.Contains(err.Error(), "expected a file") {
			t.Fatalf("unexpected error for directory path: %q", err.Error())
		}
	})
}

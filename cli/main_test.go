package main

import "testing"

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

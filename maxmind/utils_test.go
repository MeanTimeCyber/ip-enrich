package maxmind

import (
	"strings"
	"testing"
)

func TestSanitizeTerminalText(t *testing.T) {
	input := "ACME\x1b[31m Corp\nLine\tTwo\x07"
	got := SanitizeTerminalText(input)

	if strings.ContainsRune(got, '\x1b') {
		t.Fatalf("expected ANSI escape to be removed, got %q", got)
	}

	if strings.ContainsRune(got, '\n') || strings.ContainsRune(got, '\t') {
		t.Fatalf("expected control whitespace to be normalized, got %q", got)
	}

	if got != "ACME[31m Corp Line Two" {
		t.Fatalf("unexpected sanitized string: %q", got)
	}
}

func TestGetDataAsMarkdownTableEscapesCells(t *testing.T) {
	results := []Result{
		{
			Domain: "bad|domain\nname",
			IP:     "198.51.100.1",
			ASN: &ASN{
				AutonomousSystemNumber:       64500,
				AutonomousSystemOrganization: "Org|Name\x1b[31m",
			},
		},
	}

	table := GetDataAsMarkdownTable(results)

	if !strings.Contains(table, "bad\\|domain name") {
		t.Fatalf("expected escaped and sanitized domain in table, got %q", table)
	}

	if !strings.Contains(table, "AS64500 Org\\|Name[31m") {
		t.Fatalf("expected escaped and sanitized ASN organization in table, got %q", table)
	}

	if strings.ContainsRune(table, '\x1b') {
		t.Fatalf("expected no escape characters in markdown output, got %q", table)
	}
}

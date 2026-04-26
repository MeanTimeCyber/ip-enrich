package main

import (
	"net"
	"strings"
	"testing"
)

func TestFirstResolvedIPReturnsErrorForEmptyList(t *testing.T) {
	_, err := firstResolvedIP(nil)
	if err == nil {
		t.Fatal("expected error for empty resolved IP list")
	}
}

func TestFirstResolvedIPReturnsFirstAddress(t *testing.T) {
	ip, err := firstResolvedIP([]net.IP{net.ParseIP("203.0.113.10"), net.ParseIP("198.51.100.5")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ip != "203.0.113.10" {
		t.Fatalf("expected first IP address, got %q", ip)
	}
}

func TestLookupDomainRejectsMalformedDomainStrings(t *testing.T) {
	tests := []string{
		"",
		"not a domain",
		"example..com",
		"-example.com",
		"http://example.com",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := lookupDomain(nil, input)
			if err == nil {
				t.Fatalf("expected error for malformed domain %q", input)
			}

			if !strings.Contains(err.Error(), "not a valid domain") {
				t.Fatalf("expected validation error for malformed domain %q, got %q", input, err.Error())
			}
		})
	}
}

func TestLookupIPRejectsMalformedIPStrings(t *testing.T) {
	tests := []string{
		"",
		"not-an-ip",
		"300.1.1.1",
		"127.0.0.1:53",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := lookupIP(nil, "", input)
			if err == nil {
				t.Fatalf("expected error for malformed IP %q", input)
			}

			if !strings.Contains(err.Error(), "error parsing IP address") {
				t.Fatalf("expected parse error for malformed IP %q, got %q", input, err.Error())
			}
		})
	}
}

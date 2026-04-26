package main

import (
	"net"
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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const maxInputLineBytes = 1024 * 1024

// newInputScanner creates a new bufio.Scanner for the given file with a buffer size of maxInputLineBytes to handle long lines in the input files.
func newInputScanner(file *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), maxInputLineBytes)
	return scanner
}

// normalizeInputLine trims whitespace from the input line and checks if it is empty or a comment (starting with #). It returns the normalized line and a boolean indicating whether the line should be processed (true) or skipped (false).
func normalizeInputLine(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", false
	}

	return line, true
}

// requireReadableFileFromEnv checks that the specified environment variable is set to a path that points to a readable file. It returns the file path if valid, or an error if the environment variable is not set, if the path cannot be accessed, or if it points to a directory.
func requireReadableFileFromEnv(envName string) (string, error) {
	path := strings.TrimSpace(os.Getenv(envName))
	if path == "" {
		return "", fmt.Errorf("environment variable %s must be set", envName)
	}

	stat, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cannot access %s (%s): %w", envName, path, err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("%s (%s) points to a directory, expected a file", envName, path)
	}

	return path, nil
}

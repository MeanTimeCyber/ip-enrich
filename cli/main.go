package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
)

const maxInputLineBytes = 1024 * 1024

func newInputScanner(file *os.File) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), maxInputLineBytes)
	return scanner
}

func normalizeInputLine(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", false
	}

	return line, true
}

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

func main() {
	// Define a flag for the IP address to lookup
	var ipString, domainString string
	var ipList, domainList string

	var maxmindDBInfo bool
	var jsonOutput bool
	var markdownOutput bool

	// Define command-line flags for the IP address and domain to lookup, as well as options for printing MaxMind DB info and JSON output
	flag.StringVar(&ipString, "i", "", "IP address to lookup")
	flag.StringVar(&domainString, "d", "", "Domain to lookup")
	flag.StringVar(&ipList, "il", "", "File containing list of IP addresses to lookup (one per line)")
	flag.StringVar(&domainList, "dl", "", "File containing list of domains to lookup (one per line)")

	// Define flags for printing MaxMind DB metadata information and for outputting results in JSON format
	flag.BoolVar(&maxmindDBInfo, "dbinfo", false, "Print MaxMind DB metadata information")
	flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
	flag.BoolVar(&markdownOutput, "md", false, "Output results in Markdown format")

	// Parse the command-line flags
	flag.Parse()
	ipString = strings.TrimSpace(ipString)
	domainString = strings.TrimSpace(domainString)
	ipList = strings.TrimSpace(ipList)
	domainList = strings.TrimSpace(domainList)

	// Check that the user supplied an input
	if ipString == "" && domainString == "" && ipList == "" && domainList == "" {
		fmt.Println("Must supply IP address with -i, domain with -d, IP list with -il, or domain list with -dl")
		flag.Usage()
		os.Exit(1)
	}

	if jsonOutput && markdownOutput {
		fmt.Println("Only one output mode may be selected: use -json or -md")
		os.Exit(1)
	}

	cityDBPath, err := requireReadableFileFromEnv(maxmind.MaxmindCityDBEnv)
	if err != nil {
		fmt.Printf("Error validating city database path: %s\n", err)
		os.Exit(1)
	}

	asnDBPath, err := requireReadableFileFromEnv(maxmind.MaxmindASNDBEnv)
	if err != nil {
		fmt.Printf("Error validating ASN database path: %s\n", err)
		os.Exit(1)
	}

	resources, err := newLookupResources(cityDBPath, asnDBPath, maxmindDBInfo)
	if err != nil {
		fmt.Printf("Error initializing lookup resources: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := resources.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: error closing database readers: %s\n", closeErr)
		}
	}()

	if maxmindDBInfo {
		fmt.Println("MaxMind database metadata:")
		maxmind.PrintDBInfo(resources.cityDB)
		maxmind.PrintDBInfo(resources.asnDB)
		fmt.Println()
	}

	results := []maxmind.Result{}

	// Perform the lookups based on the supplied input
	if domainString != "" {
		res, err := lookupDomain(resources, domainString)

		if err != nil {
			fmt.Printf("Error looking up domain %q: %s\n", domainString, err.Error())
			os.Exit(1)
		}

		results = append(results, *res)
	} else if domainList != "" { // read file line by line and lookup each domain
		// read file line by line and lookup each domain
		file, err := os.Open(domainList)
		if err != nil {
			fmt.Printf("Error opening domain list file %q: %s\n", domainList, err.Error())
			os.Exit(1)
		}
		defer file.Close()

		scanner := newInputScanner(file)

		// read file line by line and lookup each domain
		for scanner.Scan() {
			line, ok := normalizeInputLine(scanner.Text())
			if !ok {
				continue
			}

			res, err := lookupDomain(resources, line)

			if err != nil {
				fmt.Printf("Error looking up domain %q: %s\n", line, err.Error())
				continue
			}

			results = append(results, *res)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading domain list file %q: %s\n", domainList, err.Error())
			os.Exit(1)
		}
	} else if ipString != "" { // Perform the lookup for the single supplied IP address
		res, err := lookupIP(resources, "", ipString)

		if err != nil {
			fmt.Printf("Error looking up IP address %q: %s\n", ipString, err.Error())
			os.Exit(1)
		}

		results = append(results, *res)
	} else if ipList != "" {
		// Perform the lookups for the list of IP addresses supplied in the file
		// read file line by line and lookup each IP address
		file, err := os.Open(ipList)

		if err != nil {
			fmt.Printf("Error opening IP list file %q: %s\n", ipList, err.Error())
			os.Exit(1)
		}

		defer file.Close()

		// read file line by line and lookup each IP address
		scanner := newInputScanner(file)

		for scanner.Scan() {
			line, ok := normalizeInputLine(scanner.Text())
			if !ok {
				continue
			}

			res, err := lookupIP(resources, "", line)
			if err != nil {
				fmt.Printf("Error looking up IP address %q: %s\n", line, err.Error())
				continue
			}

			results = append(results, *res)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading IP list file %q: %s\n", ipList, err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("Done with lookups. Got %d results.\n\n", len(results))

	// Output the results in the specified format
	if jsonOutput {
		outputJSON(results)
	} else if markdownOutput {
		outputMarkdown(results)
	} else {
		outputHumanReadable(results)
	}

	fmt.Printf("Fin.\n")
}

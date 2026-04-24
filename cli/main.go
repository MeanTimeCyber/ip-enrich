package main

import (
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
)

func main() {
	// Define a flag for the IP address to lookup
	var ipString, domainString string
	flag.StringVar(&ipString, "i", "", "IP address to lookup")
	flag.StringVar(&domainString, "d", "", "Domain to lookup")
	flag.Parse()

	// Check that the user supplied an IP address or a domain to lookup
	if ipString == "" && domainString == "" {
		fmt.Println("Must supply IP address with -i or domain with -d")
		flag.Usage()
		os.Exit(-1)
	}

	var ip netip.Addr

	if ipString != "" {
		fmt.Printf("Looking up address: %s\n", ipString)

		// validate the IP address format
		ip = checkAndParseAddressString(ipString)
	} else {
		// resolve the domain to an IP address
		resolvedIP, err := net.LookupIP(domainString)

		if err != nil {
			fmt.Printf("Error resolving domain %q to an address: %s\n", domainString, err.Error())
			os.Exit(-1)
		}

		// Use the first resolved IP address for the lookups
		fmt.Printf("Got address %s for domain %q\n", resolvedIP[0].String(), domainString)
		ip = checkAndParseAddressString(resolvedIP[0].String())
	}

	// Perform the MaxMind lookups and print the results
	printCity(ip)
	printASN(ip)
}

// printCity performs the MaxMind City lookup for the given IP address and prints the results.
func printCity(ip netip.Addr) {
	// Get the MaxMind City DB path from the environment variable
	dbPath := os.Getenv(maxmind.MaxmindCityDBEnv)
	if dbPath == "" {
		fmt.Printf("Must set MaxMind City DB path with the env variable, e.g. 'export %s=<path to db>'\n", maxmind.MaxmindCityDBEnv)
		os.Exit(-1)
	}

	// Perform the IP lookups
	city, err := maxmind.GetCityFromIP(ip, dbPath)
	if err != nil {
		fmt.Printf("Error looking up IP address: %v\n", err)
		os.Exit(-1)
	}

	// Print the results
	fmt.Println("\n---- Geo Lookup ----")
	maxmind.PrintCityDetails(city)
}

// printASN performs the MaxMind ASN lookup for the given IP address and prints the results.
func printASN(ip netip.Addr) {
	// Get the MaxMind ASN DB path from the environment variable
	dbPath := os.Getenv(maxmind.MaxmindASNDBEnv)
	if dbPath == "" {
		fmt.Printf("Must set MaxMind ASN DB path with the env variable, e.g. 'export %s=<path to db>'\n", maxmind.MaxmindASNDBEnv)
		os.Exit(-1)
	}

	// Perform the IP lookups
	asn, err := maxmind.GetASNFromIP(ip, dbPath)
	if err != nil {
		fmt.Printf("Error looking up IP address: %v\n", err)
		os.Exit(-1)
	}

	// Print the results
	fmt.Println("\n---- ASN Lookup ----")
	maxmind.PrintASNDetails(asn)
}

// checkAndParseAddressString validates the IP address format and returns a netip.Addr if valid,
// otherwise it prints an error message and exits the program.
func checkAndParseAddressString(ipString string) netip.Addr {
	ip, err := netip.ParseAddr(ipString)

	if err != nil {
		fmt.Printf("Invalid IP address format: %s\n", ipString)
		os.Exit(-1)
	}

	return ip
}

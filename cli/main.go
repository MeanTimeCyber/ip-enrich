package main

import (
	"flag"
	"fmt"
	"net/netip"
	"os"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
)

const (
	maxmindCityDBEnv = "MAXMIND_CITY_DB"
	maxmindASNDBEnv  = "MAXMIND_ASN_DB"
)

func main() {
	// Define a flag for the IP address to lookup
	var ipString string
	flag.StringVar(&ipString, "i", "", "IP address to lookup")
	flag.Parse()

	// Check that the user supplied an IP address
	if ipString == "" {
		fmt.Println("Must supply IP address with -i")
		flag.Usage()
		os.Exit(-1)
	}

	// validate the IP address format
	ip := checkAndParseAddressString(ipString)

	printCity(ip)
	printASN(ip)
}

// printCity performs the MaxMind City lookup for the given IP address and prints the results.
func printCity(ip netip.Addr) {
	// Get the MaxMind City DB path from the environment variable
	dbPath := os.Getenv(maxmindCityDBEnv)
	if dbPath == "" {
		fmt.Printf("Must supply MaxMind City DB path with %s\n", maxmindCityDBEnv)
		os.Exit(-1)
	}

	// Perform the IP lookups
	city, err := maxmind.GetCityFromIP(ip, dbPath)
	if err != nil {
		fmt.Printf("Error looking up IP address: %v\n", err)
		os.Exit(-1)
	}

	// Print the results
	fmt.Println("\n---- Geo-lookup ----")
	maxmind.PrintCityDetails(city)
}

// printASN performs the MaxMind ASN lookup for the given IP address and prints the results.
func printASN(ip netip.Addr) {
	// Get the MaxMind ASN DB path from the environment variable
	dbPath := os.Getenv(maxmindASNDBEnv)
	if dbPath == "" {
		fmt.Printf("Must supply MaxMind ASN DB path with %s\n", maxmindASNDBEnv)
		os.Exit(-1)
	}

	// Perform the IP lookups
	asn, err := maxmind.GetASNFromIP(ip, dbPath)
	if err != nil {
		fmt.Printf("Error looking up IP address: %v\n", err)
		os.Exit(-1)
	}

	// Print the results
	fmt.Println("\n---- ASN ----")
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

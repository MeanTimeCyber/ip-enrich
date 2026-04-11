package main

import (
	"flag"
	"fmt"
	"net/netip"
	"os"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
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

	// Get the MaxMind City DB path from the environment variable
	dbPath := os.Getenv("MAXMINDCITYDB")
	if dbPath == "" {
		fmt.Println("Must supply MaxMind City DB path with MAXMINDCITYDB")
		os.Exit(-1)
	}

	// Perform the IP lookups
	city, err := maxmind.GetCityFromIP(ip, dbPath)
	if err != nil {
		fmt.Printf("Error looking up IP address: %v\n", err)
		os.Exit(-1)
	}

	// Print the results
	maxmind.PrintCityDetails(city)
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

package main

import (
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
	"github.com/asaskevich/govalidator"
)

func main() {
	// Define a flag for the IP address to lookup
	var ipString, domainString string
	var ipList, domainList string
	var maxmindDBInfo bool
	var jsonOutput bool

	// Define command-line flags for the IP address and domain to lookup, as well as options for printing MaxMind DB info and JSON output
	flag.StringVar(&ipString, "i", "", "IP address to lookup")
	flag.StringVar(&domainString, "d", "", "Domain to lookup")
	flag.StringVar(&ipList, "il", "", "File containing list of IP addresses to lookup (one per line)")
	flag.StringVar(&domainList, "dl", "", "File containing list of domains to lookup (one per line)")

	// Define flags for printing MaxMind DB metadata information and for outputting results in JSON format
	flag.BoolVar(&maxmindDBInfo, "dbinfo", false, "Print MaxMind DB metadata information")
	flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")

	// Parse the command-line flags
	flag.Parse()

	// Check that the user supplied an input
	if ipString == "" && domainString == "" && ipList == "" && domainList == "" {
		fmt.Println("Must supply IP address with -i, domain with -d, IP list with -il, or domain list with -dl")
		flag.Usage()
		os.Exit(-1)
	}

	if domainString != "" {
		lookupDomain(domainString, maxmindDBInfo, jsonOutput)
	}

	if ipString != "" {
		lookupIP(ipString, maxmindDBInfo, jsonOutput)
	}

	fmt.Printf("\nFin.\n")
}

func lookupDomain(domainString string, maxmindDBInfo, jsonOutput bool) {
	// check domain validity if supplied
	if !govalidator.IsDNSName(domainString) {
		fmt.Printf("%s is not a valid domain\n", domainString)
		os.Exit(-1)
	}

	// resolve the domain to an IP address
	resolvedIP, err := net.LookupIP(domainString)

	if err != nil {
		fmt.Printf("Error resolving domain %q to an address: %s\n", domainString, err.Error())
		os.Exit(-1)
	}

	// Use the first resolved IP address for the lookups
	fmt.Printf("Got address %s for domain %q\n", resolvedIP[0].String(), domainString)
	lookupIP(resolvedIP[0].String(), maxmindDBInfo, jsonOutput)
}

func lookupIP(ipString string, maxmindDBInfo, jsonOutput bool) {

	fmt.Printf("Looking up address: %s\n", ipString)

	// validate the IP address format
	ip := checkAndParseAddressString(ipString)

	// Perform the MaxMind lookups for the IP address
	city, err := maxmind.GetCityFromIP(ip, os.Getenv(maxmind.MaxmindCityDBEnv), maxmindDBInfo)

	if err != nil {
		fmt.Printf("Error looking up city for IP address: %v\n", err)
		os.Exit(-1)
	}

	// Perform the MaxMind ASN lookup for the IP address
	asn, err := maxmind.GetASNFromIP(ip, os.Getenv(maxmind.MaxmindASNDBEnv), maxmindDBInfo)

	if err != nil {
		fmt.Printf("Error looking up ASN for IP address: %v\n", err)
		os.Exit(-1)
	}

	// Output the results in JSON format if the flag is set, otherwise print them in a human-readable format
	if jsonOutput {
		jsonString, err := maxmind.GetDataAsFormattedJSON(city, asn)

		if err != nil {
			fmt.Printf("Error generating JSON output: %v\n", err)
			os.Exit(-1)
		}

		fmt.Println(jsonString)
	} else {
		// Print the results to the console in a human-readable format
		fmt.Println("\n---- Geo Lookup ----")
		maxmind.PrintCityDetails(city)

		fmt.Println("\n---- ASN Lookup ----")
		maxmind.PrintASNDetails(asn)
	}
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

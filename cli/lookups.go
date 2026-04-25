package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
	"github.com/asaskevich/govalidator"
)

// lookupDomain resolves the given domain to an IP address and performs the MaxMind lookups for that IP address.
func lookupDomain(domainString string, maxmindDBInfo, jsonOutput bool) error {
	// check domain validity if supplied
	if !govalidator.IsDNSName(domainString) {
		return fmt.Errorf("%s is not a valid domain: %s\n", domainString)
	}

	fmt.Printf("Resolving domain %q\n", domainString)

	// resolve the domain to an IP address
	resolvedIP, err := net.LookupIP(domainString)

	if err != nil {
		return fmt.Errorf("Error resolving domain %q to an address: %s", domainString, err.Error())
	}

	// Use the first resolved IP address for the lookups
	fmt.Printf("Got address %s for domain %q\n\n", resolvedIP[0].String(), domainString)
	err = lookupIP(resolvedIP[0].String(), maxmindDBInfo, jsonOutput)

	return err
}

// lookupIP performs the MaxMind lookups for the given IP address and outputs the results in either JSON format or a human-readable format.
func lookupIP(ipString string, maxmindDBInfo, jsonOutput bool) error {
	fmt.Printf("Looking up address: %s\n", ipString)

	// validate the IP address format
	ip, err := netip.ParseAddr(ipString)

	if err != nil {
		return fmt.Errorf("Error parsing IP address %q: %s", ipString, err.Error())
	}

	fmt.Printf("Resolving domain %q\n", ipString)

	// Perform the MaxMind lookups for the IP address
	city, err := maxmind.GetCityFromIP(ip, os.Getenv(maxmind.MaxmindCityDBEnv), maxmindDBInfo)

	if err != nil {
		return fmt.Errorf("Error looking up city for IP address %q: %v", ipString, err)
	}

	// Perform the MaxMind ASN lookup for the IP address
	asn, err := maxmind.GetASNFromIP(ip, os.Getenv(maxmind.MaxmindASNDBEnv), maxmindDBInfo)

	if err != nil {
		return fmt.Errorf("Error looking up ASN for IP address %q: %v", ipString, err)
	}

	// Output the results in JSON format if the flag is set, otherwise print them in a human-readable format
	if jsonOutput {
		jsonString, err := maxmind.GetDataAsFormattedJSON(city, asn)

		if err != nil {
			return fmt.Errorf("Error generating JSON output: %v", err)
		}

		fmt.Println(jsonString)
	} else {
		// Print the results to the console in a human-readable format
		fmt.Println("\n---- Geo Lookup ----")
		maxmind.PrintCityDetails(city)

		fmt.Println("\n---- ASN Lookup ----")
		maxmind.PrintASNDetails(asn)
	}

	return nil
}

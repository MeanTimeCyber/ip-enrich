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
func lookupDomain(domainString string, maxmindDBInfo bool) (*maxmind.Result, error) {
	// check domain validity if supplied
	if !govalidator.IsDNSName(domainString) {
		return nil, fmt.Errorf("%s is not a valid domain\n", domainString)
	}

	fmt.Printf("Resolving domain %q\n", domainString)

	// resolve the domain to an IP address
	resolvedIP, err := net.LookupIP(domainString)

	if err != nil {
		return nil, fmt.Errorf("Error resolving domain %q to an address: %s", domainString, err.Error())
	}

	// Use the first resolved IP address for the lookups
	fmt.Printf("Got address %s for domain %q\n", resolvedIP[0].String(), domainString)

	return lookupIP(domainString, resolvedIP[0].String(), maxmindDBInfo)
}

// lookupIP performs the MaxMind lookups for the given IP address and outputs the results in either JSON format or a human-readable format.
func lookupIP(domainString, ipString string, maxmindDBInfo bool) (*maxmind.Result, error) {
	fmt.Printf("Looking up address: %s\n", ipString)

	// validate the IP address format
	ip, err := netip.ParseAddr(ipString)

	if err != nil {
		return nil, fmt.Errorf("Error parsing IP address %q: %s", ipString, err.Error())
	}

	fmt.Printf("Resolving IP %q\n", ipString)

	// Perform the MaxMind lookups for the IP address
	city, err := maxmind.GetCityFromIP(ip, os.Getenv(maxmind.MaxmindCityDBEnv), maxmindDBInfo)

	if err != nil {
		return nil, fmt.Errorf("Error looking up city for IP address %q: %v", ipString, err)
	}

	// Perform the MaxMind ASN lookup for the IP address
	asn, err := maxmind.GetASNFromIP(ip, os.Getenv(maxmind.MaxmindASNDBEnv), maxmindDBInfo)

	if err != nil {
		return nil, fmt.Errorf("Error looking up ASN for IP address %q: %v", ipString, err)
	}

	fmt.Printf("Got city and ASN information for IP %q\n\n", ipString)

	// Create a Result struct to hold the lookup results and return it
	return &maxmind.Result{
		Domain: domainString,
		IP:     ipString,
		City:   city,
		ASN:    asn,
	}, nil
}

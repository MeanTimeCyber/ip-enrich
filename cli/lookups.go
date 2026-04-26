package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
	"github.com/asaskevich/govalidator"
	"github.com/oschwald/maxminddb-golang/v2"
)

const defaultDNSTimeout = 5 * time.Second

type lookupResources struct {
	cityDB     *maxminddb.Reader
	asnDB      *maxminddb.Reader
	dnsTimeout time.Duration
	resolver   *net.Resolver
	printInfo  bool
}

func newLookupResources(cityDBPath, asnDBPath string, printInfo bool) (*lookupResources, error) {
	cityDB, err := maxminddb.Open(cityDBPath)
	if err != nil {
		return nil, fmt.Errorf("error opening city database: %w", err)
	}

	asnDB, err := maxminddb.Open(asnDBPath)
	if err != nil {
		cityDB.Close()
		return nil, fmt.Errorf("error opening ASN database: %w", err)
	}

	return &lookupResources{
		cityDB:     cityDB,
		asnDB:      asnDB,
		dnsTimeout: defaultDNSTimeout,
		resolver:   net.DefaultResolver,
		printInfo:  printInfo,
	}, nil
}

func (resources *lookupResources) Close() error {
	if resources == nil {
		return nil
	}

	var errs []error
	if resources.cityDB != nil {
		if err := resources.cityDB.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if resources.asnDB != nil {
		if err := resources.asnDB.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func firstResolvedIP(resolvedIP []net.IP) (string, error) {
	if len(resolvedIP) == 0 {
		return "", fmt.Errorf("no addresses resolved")
	}

	return resolvedIP[0].String(), nil
}

// lookupDomain resolves the given domain to an IP address and performs the MaxMind lookups for that IP address.
func lookupDomain(resources *lookupResources, domainString string) (*maxmind.Result, error) {
	// check domain validity if supplied
	if !govalidator.IsDNSName(domainString) {
		return nil, fmt.Errorf("%q is not a valid domain", domainString)
	}

	fmt.Printf("Resolving domain %q\n", maxmind.SanitizeTerminalText(domainString))

	ctx, cancel := context.WithTimeout(context.Background(), resources.dnsTimeout)
	defer cancel()

	resolvedIP, err := resources.resolver.LookupIP(ctx, "ip", domainString)

	if err != nil {
		return nil, fmt.Errorf("error resolving domain %q to an address: %w", domainString, err)
	}

	lookupAddr, err := firstResolvedIP(resolvedIP)
	if err != nil {
		return nil, fmt.Errorf("error resolving domain %q to an address: %w", domainString, err)
	}

	fmt.Printf("Got address %s for domain %q\n", lookupAddr, maxmind.SanitizeTerminalText(domainString))

	return lookupIP(resources, domainString, lookupAddr)
}

// lookupIP performs the MaxMind lookups for the given IP address and outputs the results in either JSON format or a human-readable format.
func lookupIP(resources *lookupResources, domainString, ipString string) (*maxmind.Result, error) {
	fmt.Printf("Looking up address: %s\n", ipString)

	// validate the IP address format
	ip, err := netip.ParseAddr(ipString)

	if err != nil {
		return nil, fmt.Errorf("error parsing IP address %q: %w", ipString, err)
	}

	fmt.Printf("Resolving IP %q\n", ipString)

	// Perform the MaxMind lookups for the IP address
	city, err := maxmind.GetCityFromIPWithReader(ip, resources.cityDB, resources.printInfo)

	if err != nil {
		return nil, fmt.Errorf("error looking up city for IP address %q: %w", ipString, err)
	}

	// Perform the MaxMind ASN lookup for the IP address
	asn, err := maxmind.GetASNFromIPWithReader(ip, resources.asnDB, resources.printInfo)

	if err != nil {
		return nil, fmt.Errorf("error looking up ASN for IP address %q: %w", ipString, err)
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

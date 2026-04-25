package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

	// Perform the lookups based on the supplied input
	if domainString != "" {
		if err := lookupDomain(domainString, maxmindDBInfo, jsonOutput); err != nil {
			fmt.Printf("Error looking up domain %q: %s\n", domainString, err.Error())
			os.Exit(-1)
		}
	} else if domainList != "" { // read file line by line and lookup each domain
		// read file line by line and lookup each domain
		file, err := os.Open(domainList)
		if err != nil {
			fmt.Printf("Error opening domain list file %q: %s\n", domainList, err.Error())
			os.Exit(-1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		
		// read file line by line and lookup each domain
		for scanner.Scan() {
			line := scanner.Text()
			
			if line != "" {
				if err := lookupDomain(line, maxmindDBInfo, jsonOutput); err != nil {
					fmt.Printf("Error looking up domain %q: %s\n", line, err.Error())
					os.Exit(-1)
				}
				fmt.Println()
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading domain list file %q: %s\n", domainList, err.Error())
			os.Exit(-1)
		}
	} else if ipString != "" { // Perform the lookup for the single supplied IP address
		if err := lookupIP(ipString, maxmindDBInfo, jsonOutput); err != nil {
			fmt.Printf("Error looking up IP address %q: %s\n", ipString, err.Error())
			os.Exit(-1)
		}

		fmt.Println()
	} else if ipList != "" { // Perform the lookups for the list of IP addresses supplied in the file
		// read file line by line and lookup each IP address
		file, err := os.Open(ipList)
		
		if err != nil {
			fmt.Printf("Error opening IP list file %q: %s\n", ipList, err.Error())
			os.Exit(-1)
		}
		
		defer file.Close()

		// read file line by line and lookup each IP address
		scanner := bufio.NewScanner(file)
		
		for scanner.Scan() {
			line := scanner.Text()
			
			if line != "" {
				if err := lookupIP(line, maxmindDBInfo, jsonOutput); err != nil {
					fmt.Printf("Error looking up IP address %q: %s\n", line, err.Error())
					os.Exit(-1)
				}
				
				fmt.Println()
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading IP list file %q: %s\n", ipList, err.Error())
			os.Exit(-1)
		}
	}

	fmt.Printf("Fin.\n")
}

package main

import (
	"fmt"
	"os"

	"github.com/MeanTimeCyber/ip-enrich/maxmind"
	"github.com/markkurossi/tabulate"
)

// outputJSON outputs the MaxMind lookup results in a JSON format.
func outputJSON(results []maxmind.Result) {
	jsonString, err := maxmind.GetDataAsFormattedJSON(results)

	if err != nil {
		fmt.Printf("Error formatting results as JSON: %s\n", err.Error())
		return
	}

	fmt.Println(jsonString)
}

// outputFormattedJSON outputs the MaxMind lookup results in a pretty-printed JSON format.
func outputMarkdown(results []maxmind.Result) {
	mdTable := maxmind.GetDataAsMarkdownTable(results)
	fmt.Println(mdTable)
}

// outputHumanReadable outputs the MaxMind lookup results in a human-readable format, which is a table with columns for the domain, IP address, country, city, subdivision, and ASN.
func outputHumanReadable(results []maxmind.Result) {
	// Print a table for the city information
	table := tabulate.New(tabulate.Unicode)
	table.Header("IP")
	table.Header("Country")
	table.Header("Subdivision")
	table.Header("City")
	table.Header("Name")
	table.Header("Postal Code")

	for _, result := range results {
		row := table.Row()
		row.Column(maxmind.SanitizeTerminalText(result.IP))
		row.Column(resultCountry(result))
		row.Column(resultSubdivision(result))
		row.Column(resultCity(result))
		row.Column(resultName(result))
		row.Column(resultPostalCode(result))
	}

	table.Print(os.Stdout)
	fmt.Println()

	// Print a separate table for ASN information	
	table = tabulate.New(tabulate.Unicode)
	table.Header("IP")
	table.Header("ASN")

	for _, result := range results {
		row := table.Row()
		row.Column(maxmind.SanitizeTerminalText(result.IP))
		row.Column(resultASN(result))
	}

	table.Print(os.Stdout)
	fmt.Println()

	// TODO something with reverse domain lookups, if we have them
}

func resultCountry(result maxmind.Result) string {
	if result.City == nil {
		return ""
	}

	return maxmind.SanitizeTerminalText(maxmind.EnglishName(result.City.Country.Names))
}

func resultCity(result maxmind.Result) string {
	if result.City == nil {
		return ""
	}

	return maxmind.SanitizeTerminalText(maxmind.DisplayCityName(*result.City))
}

func resultSubdivision(result maxmind.Result) string {
	if result.City == nil || len(result.City.Subdivisions) == 0 {
		return ""
	}

	return maxmind.SanitizeTerminalText(maxmind.GetSubdivisionValue(result.City.Subdivisions[0]))
}

func resultASN(result maxmind.Result) string {
	if result.ASN == nil {
		return ""
	}

	return maxmind.SanitizeTerminalText(fmt.Sprintf("AS%d %s", result.ASN.AutonomousSystemNumber, result.ASN.AutonomousSystemOrganization))
}

func resultName(result maxmind.Result) string {
	if result.City == nil {
		return ""
	}

	return maxmind.SanitizeTerminalText(maxmind.EnglishName(result.City.City.Names))
}

func resultPostalCode(result maxmind.Result) string {
	if result.City == nil {
		return ""
	}

	return maxmind.SanitizeTerminalText(result.City.Postal.Code)
}

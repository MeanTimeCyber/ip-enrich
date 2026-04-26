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
	table := tabulate.New(tabulate.Simple)
	table.Header("Domain")
	table.Header("IP")
	table.Header("Country")
	table.Header("City")
	table.Header("Subdivision")
	table.Header("ASN")

	for _, result := range results {
		row := table.Row()
		row.Column(result.Domain)
		row.Column(result.IP)
		row.Column(resultCountry(result))
		row.Column(resultCity(result))
		row.Column(resultSubdivision(result))
		row.Column(resultASN(result))
	}

	table.Print(os.Stdout)
	fmt.Println()
}

func resultCountry(result maxmind.Result) string {
	if result.City == nil {
		return ""
	}

	return maxmind.EnglishName(result.City.Country.Names)
}

func resultCity(result maxmind.Result) string {
	if result.City == nil {
		return ""
	}

	return maxmind.DisplayCityName(*result.City)
}

func resultSubdivision(result maxmind.Result) string {
	if result.City == nil || len(result.City.Subdivisions) == 0 {
		return ""
	}

	return maxmind.GetSubdivisionValue(result.City.Subdivisions[0])
}

func resultASN(result maxmind.Result) string {
	if result.ASN == nil {
		return ""
	}

	return fmt.Sprintf("AS%d %s", result.ASN.AutonomousSystemNumber, result.ASN.AutonomousSystemOrganization)
}

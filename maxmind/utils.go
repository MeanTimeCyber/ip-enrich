package maxmind

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

// printDBInfo prints the metadata information of the MaxMind database, including the database type, IP version, and build time. It also checks if the database is more than 3 months old and prints a warning if it is.
func printDBInfo(db *maxminddb.Reader) {
	meta := db.Metadata
	buildTime := meta.BuildTime()
	fmt.Printf("Database: %s (IPv%d, built %s)\n", meta.DatabaseType, meta.IPVersion, buildTime.Format("2006-01-02"))

	if buildTime.Before(time.Now().AddDate(0, -3, 0)) {
		fmt.Fprintln(os.Stderr, "Warning: database is more than 3 months old")
	}
}

// EnglishName returns the English name from the given map of names, or an empty string if there is no English name or if the map is empty.
func EnglishName(names map[string]string) string {
	if len(names) == 0 {
		return ""
	}

	if name := names["en"]; name != "" {
		return name
	}

	for _, name := range names {
		if name != "" {
			return name
		}
	}

	return ""
}

type tableRow struct {
	label string
	value string
}

// postalArea extracts the postal area from the given postal code, which is the leading letters of the postal code before any digits.
func postalArea(postalCode string) string {
	var builder strings.Builder
	for _, char := range postalCode {
		if char >= '0' && char <= '9' {
			break
		}
		if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' {
			builder.WriteRune(char)
		}
	}

	return strings.ToUpper(builder.String())
}

// isLondonPostalArea checks if the given postal code is in a London postal area, which includes the postal areas E, EC, N, NW, SE, SW, W, and WC.
func isLondonPostalArea(postalCode string) bool {
	switch postalArea(postalCode) {
	case "E", "EC", "N", "NW", "SE", "SW", "W", "WC":
		return true
	default:
		return false
	}
}

// DisplayCityName returns the city name to display, which is the city name if it is available,
// otherwise it falls back to the locality name, and if that is not available it returns an empty string.
func DisplayCityName(record City) string {
	if record.Country.ISOCode == "GB" && isLondonPostalArea(record.Postal.Code) {
		return "London"
	}

	return EnglishName(record.City.Names)
}

// GetSubdivisionValue returns the subdivision name to display, which is the subdivision name if it is available,
func GetSubdivisionValue(subdivision struct {
	GeoNameID uint32            `maxminddb:"geoname_id"`
	ISOCode   string            `maxminddb:"iso_code"`
	Names     map[string]string `maxminddb:"names"`
}) string {
	name := EnglishName(subdivision.Names)
	if subdivision.ISOCode == "" {
		return name
	}

	return fmt.Sprintf("%s (%s)", name, subdivision.ISOCode)
}

// GetDataAsFormattedJSON returns the MaxMind lookup results formatted as a JSON string.
func GetDataAsJSON(results []Result) (string, error) {
	jsonBytes, err := json.Marshal(results)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// GetDataAsFormattedJSON returns the MaxMind lookup results formatted as a pretty-printed JSON string.
func GetDataAsFormattedJSON(results []Result) (string, error) {
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// GetDataAsMarkdownTable returns the MaxMind lookup results formatted as a Markdown table string.
func GetDataAsMarkdownTable(results []Result) string {
	var builder strings.Builder

	builder.WriteString("| Domain | IP | Country | City | Subdivision | ASN |\n")
	builder.WriteString("|--------|----|---------|------|-------------|-----|\n")

	for _, result := range results {
		builder.WriteString("|")
		builder.WriteString(result.Domain)
		builder.WriteString("|")
		builder.WriteString(result.IP)
		builder.WriteString("|")

		if result.City != nil {
			builder.WriteString(EnglishName(result.City.Country.Names))
			builder.WriteString(" | ")
			builder.WriteString(DisplayCityName(*result.City))
			builder.WriteString(" | ")
			if len(result.City.Subdivisions) > 0 {
				builder.WriteString(GetSubdivisionValue(result.City.Subdivisions[0]))
			} else {
				builder.WriteString("")
			}
		} else {
			builder.WriteString(" |  | ")
		}

		if result.ASN != nil {
			builder.WriteString(fmt.Sprintf("AS%d %s", result.ASN.AutonomousSystemNumber, result.ASN.AutonomousSystemOrganization))
		} else {
			builder.WriteString("")
		}

		builder.WriteString("|\n")
	}

	return builder.String()
}

package maxmind

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/oschwald/maxminddb-golang/v2"
)

func printDBInfo(db *maxminddb.Reader) {
	meta := db.Metadata
	buildTime := meta.BuildTime()
	//fmt.Printf("Database: %s (IPv%d, built %s)\n", meta.DatabaseType, meta.IPVersion, buildTime.Format("2006-01-02"))
	if buildTime.Before(time.Now().AddDate(0, -3, 0)) {
		fmt.Fprintln(os.Stderr, "Warning: database is more than 3 months old")
	}
}

func englishName(names map[string]string) string {
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

func isLondonPostalArea(postalCode string) bool {
	switch postalArea(postalCode) {
	case "E", "EC", "N", "NW", "SE", "SW", "W", "WC":
		return true
	default:
		return false
	}
}

func displayCityName(record City) string {
	if record.Country.ISOCode == "GB" && isLondonPostalArea(record.Postal.Code) {
		return "London"
	}

	return englishName(record.City.Names)
}

func subdivisionValue(subdivision struct {
	GeoNameID uint32            `maxminddb:"geoname_id"`
	ISOCode   string            `maxminddb:"iso_code"`
	Names     map[string]string `maxminddb:"names"`
}) string {
	name := englishName(subdivision.Names)
	if subdivision.ISOCode == "" {
		return name
	}

	return fmt.Sprintf("%s (%s)", name, subdivision.ISOCode)
}

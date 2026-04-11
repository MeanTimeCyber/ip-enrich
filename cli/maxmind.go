package main

import (
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/markkurossi/tabulate"
	"github.com/oschwald/maxminddb-golang/v2"
)

type City struct {
	City struct {
		GeoNameID uint32            `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Continent struct {
		Code      string            `maxminddb:"code"`
		GeoNameID uint32            `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		GeoNameID           uint32            `maxminddb:"geoname_id"`
		ISOCode             string            `maxminddb:"iso_code"`
		IsInEuropeanUnion   bool              `maxminddb:"is_in_european_union"`
		Names               map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		TimeZone       string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	RegisteredCountry struct {
		GeoNameID         uint32            `maxminddb:"geoname_id"`
		ISOCode           string            `maxminddb:"iso_code"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		Names             map[string]string `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID         uint32            `maxminddb:"geoname_id"`
		ISOCode           string            `maxminddb:"iso_code"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		Names             map[string]string `maxminddb:"names"`
		Type              string            `maxminddb:"type"`
	} `maxminddb:"represented_country"`
	Subdivisions []struct {
		GeoNameID uint32            `maxminddb:"geoname_id"`
		ISOCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
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

// https://www.maxmind.com/en/geoip-databases
func lookupMaxMindCity(ip netip.Addr, dbPath string) error {
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	var city City
	err = db.Lookup(ip).Decode(&city)

	if err != nil {
		return err
	}

	displayCity := displayCityName(city)
	rawLocality := englishName(city.City.Names)

	rows := []tableRow{
		{label: "City", value: displayCity},
		{label: "Country", value: fmt.Sprintf("%s (%s)", englishName(city.Country.Names), city.Country.ISOCode)},
	}

	if rawLocality != "" && rawLocality != displayCity {
		rows = append(rows, tableRow{label: "Locality", value: rawLocality})
	}

	if len(city.Subdivisions) > 0 {
		rows = append(rows, tableRow{label: "Region", value: subdivisionValue(city.Subdivisions[0])})
	}

	if len(city.Subdivisions) > 1 {
		rows = append(rows, tableRow{label: "District", value: subdivisionValue(city.Subdivisions[1])})
	}

	for index := 2; index < len(city.Subdivisions); index++ {
		rows = append(rows, tableRow{label: fmt.Sprintf("Subdivision %d", index+1), value: subdivisionValue(city.Subdivisions[index])})
	}

	if city.Continent.Code != "" || len(city.Continent.Names) > 0 {
		rows = append(rows, tableRow{label: "Continent", value: fmt.Sprintf("%s (%s)", englishName(city.Continent.Names), city.Continent.Code)})
	}

	if city.Postal.Code != "" {
		rows = append(rows, tableRow{label: "Postal Code", value: city.Postal.Code})
	}

	if city.Location.TimeZone != "" {
		rows = append(rows, tableRow{label: "Time Zone", value: city.Location.TimeZone})
	}

	if city.Location.Latitude != 0 || city.Location.Longitude != 0 {
		rows = append(rows, tableRow{label: "Coordinates", value: fmt.Sprintf("%.4f, %.4f", city.Location.Latitude, city.Location.Longitude)})
	}

	if city.Location.AccuracyRadius != 0 {
		rows = append(rows, tableRow{label: "Accuracy Radius", value: fmt.Sprintf("%d km", city.Location.AccuracyRadius)})
	}

	if city.RegisteredCountry.ISOCode != "" {
		rows = append(rows, tableRow{label: "Registered Country", value: fmt.Sprintf("%s (%s)", englishName(city.RegisteredCountry.Names), city.RegisteredCountry.ISOCode)})
	}

	if city.RepresentedCountry.ISOCode != "" {
		representedCountry := fmt.Sprintf("%s (%s)", englishName(city.RepresentedCountry.Names), city.RepresentedCountry.ISOCode)
		if city.RepresentedCountry.Type != "" {
			representedCountry += fmt.Sprintf(" [%s]", city.RepresentedCountry.Type)
		}
		rows = append(rows, tableRow{label: "Represented Country", value: representedCountry})
	}

	if city.Country.IsInEuropeanUnion {
		rows = append(rows, tableRow{label: "In European Union", value: "yes"})
	}

	table := tabulate.New(tabulate.Simple)
	table.Header("Field").SetAlign(tabulate.ML)
	table.Header("Value").SetAlign(tabulate.ML)

	for _, row := range rows {
		dataRow := table.Row()
		dataRow.Column(row.label)
		dataRow.Column(row.value)
	}

	table.Print(os.Stdout)
	return nil
}

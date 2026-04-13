package maxmind

import (
	"fmt"
	"net/netip"
	"os"

	"github.com/markkurossi/tabulate"
	"github.com/oschwald/maxminddb-golang/v2"
)

// https://www.maxmind.com/en/geoip-databases
func GetCityFromIP(ip netip.Addr, dbPath string) (*City, error) {
	// Open the MaxMind database. The Open function returns a Reader
	// that can be used to perform lookups on the database.
	// The Reader must be closed when it is no longer needed.
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	printDBInfo(db)

	// Decode the record into a City struct.
	// The City struct must have fields that match the structure of the database record,
	// and the maxminddb tags must be used to specify the field names in the database.
	var city City
	err = db.Lookup(ip).Decode(&city)

	if err != nil {
		return nil, err
	}

	return &city, nil
}

// displayCityName returns the city name to display, which is the city name if it is available,
// otherwise it falls back to the locality name, and if that is not available it returns an empty string.
func PrintCityDetails(city *City) {
	displayCity := displayCityName(*city)
	rawLocality := englishName(city.City.Names)

	var rows []tableRow

	if city.Location.Latitude != 0 || city.Location.Longitude != 0 {
		rows = append(rows, tableRow{label: "Coordinates", value: fmt.Sprintf("%.4f, %.4f", city.Location.Latitude, city.Location.Longitude)})
	}

	if city.Location.AccuracyRadius != 0 {
		rows = append(rows, tableRow{label: "Accuracy Radius", value: fmt.Sprintf("%d km", city.Location.AccuracyRadius)})
	}

	if rawLocality != "" && rawLocality != displayCity {
		rows = append(rows, tableRow{label: "Locality", value: rawLocality})
	}

	if len(city.Subdivisions) > 1 {
		rows = append(rows, tableRow{label: "District", value: subdivisionValue(city.Subdivisions[1])})
	}

	for index := 2; index < len(city.Subdivisions); index++ {
		rows = append(rows, tableRow{label: fmt.Sprintf("Subdivision %d", index+1), value: subdivisionValue(city.Subdivisions[index])})
	}

	if city.Postal.Code != "" {
		rows = append(rows, tableRow{label: "Postal Code", value: city.Postal.Code})
	}

	rows = append(rows, tableRow{label: "City", value: displayCity})

	if len(city.Subdivisions) > 0 {
		rows = append(rows, tableRow{label: "Region", value: subdivisionValue(city.Subdivisions[0])})
	}

	rows = append(rows, tableRow{label: "Country", value: fmt.Sprintf("%s (%s)", englishName(city.Country.Names), city.Country.ISOCode)})

	if city.Country.IsInEuropeanUnion {
		rows = append(rows, tableRow{label: "In European Union", value: "yes"})
	}

	// if city.RegisteredCountry.ISOCode != "" {
	// 	rows = append(rows, tableRow{label: "Registered Country", value: fmt.Sprintf("%s (%s)", englishName(city.RegisteredCountry.Names), city.RegisteredCountry.ISOCode)})
	// }

	if city.RepresentedCountry.ISOCode != "" {
		representedCountry := fmt.Sprintf("%s (%s)", englishName(city.RepresentedCountry.Names), city.RepresentedCountry.ISOCode)
		if city.RepresentedCountry.Type != "" {
			representedCountry += fmt.Sprintf(" [%s]", city.RepresentedCountry.Type)
		}
		rows = append(rows, tableRow{label: "Represented Country", value: representedCountry})
	}

	if city.Continent.Code != "" || len(city.Continent.Names) > 0 {
		rows = append(rows, tableRow{label: "Continent", value: fmt.Sprintf("%s (%s)", englishName(city.Continent.Names), city.Continent.Code)})
	}

	if city.Location.TimeZone != "" {
		rows = append(rows, tableRow{label: "Time Zone", value: city.Location.TimeZone})
	}

	table := tabulate.New(tabulate.Simple)

	for _, row := range rows {
		dataRow := table.Row()
		dataRow.Column(row.label)
		dataRow.Column(row.value)
	}

	table.Print(os.Stdout)
}

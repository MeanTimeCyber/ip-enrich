package maxmind

import (
	"fmt"
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

func GetASNFromIP(ip netip.Addr, dbPath string, printInfo bool) (*ASN, error) {
	// Open the MaxMind database. The Open function returns a Reader
	// that can be used to perform lookups on the database.
	// The Reader must be closed when it is no longer needed.
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if printInfo {
		printDBInfo(db)
	}

	// Decode the record into an ASN struct.
	// The ASN struct must have fields that match the structure of the database record,
	// and the maxminddb tags must be used to specify the field names in the database.
	var asn ASN
	err = db.Lookup(ip).Decode(&asn)

	if err != nil {
		return nil, err
	}

	return &asn, nil
}

func PrintASNDetails(asn *ASN) {
	fmt.Printf("AS%d %s\n", asn.AutonomousSystemNumber, asn.AutonomousSystemOrganization)
}

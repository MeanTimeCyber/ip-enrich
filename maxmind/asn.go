package maxmind

import (
	"fmt"
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

// GetASNFromIPWithReader decodes ASN data for an IP using an already-open MaxMind reader.
func GetASNFromIPWithReader(ip netip.Addr, db *maxminddb.Reader, printInfo bool) (*ASN, error) {
	if db == nil {
		return nil, fmt.Errorf("asn database reader is nil")
	}

	if printInfo {
		printDBInfo(db)
	}

	var asn ASN
	err := db.Lookup(ip).Decode(&asn)
	if err != nil {
		return nil, err
	}

	return &asn, nil
}

func GetASNFromIP(ip netip.Addr, dbPath string, printInfo bool) (*ASN, error) {
	// Open the MaxMind database. The Open function returns a Reader
	// that can be used to perform lookups on the database.
	// The Reader must be closed when it is no longer needed.
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return GetASNFromIPWithReader(ip, db, printInfo)
}

func PrintASNDetails(asn *ASN) {
	if asn == nil {
		return
	}

	fmt.Printf("AS%d %s\n", asn.AutonomousSystemNumber, SanitizeTerminalText(asn.AutonomousSystemOrganization))
}

# ip-enrich

`ip-enrich` is a small Go CLI for enriching an IP address with location and ASN data from MaxMind databases.

## Requirements

- Go 1.25.8 or newer
- A MaxMind GeoLite2 City database file
- A MaxMind GeoLite2 ASN database file

## Build

Build the CLI from the repository root with the provided Makefile:

```bash
make build
```

This builds a small stripped binary using Go linker flags.

If you want to run the build command directly:

```bash
go build -trimpath -ldflags="-s -w" -o ip-enrich ./cli
```

This creates a local executable named `ip-enrich`.

## Run

Set the MaxMind database paths with `MAXMIND_CITY_DB` and `MAXMIND_ASN_DB`, then pass either an IP with `-i` or a domain with `-d`.

Run directly with Go:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb go run ./cli -i 8.8.8.8
```

Run the built binary:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -i 8.8.8.8
```

Look up a domain:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -d example.com
```

## Usage

```text
./ip-enrich [-i <ip-address> | -d <domain>] [-json] [-dbinfo]
```

Example:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -i 8.8.8.8
```

JSON output example:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -d example.com -json
```

Print DB metadata info while looking up an IP:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -i 1.1.1.1 -dbinfo
```

## Notes

- The City database path is required through `MAXMIND_CITY_DB`.
- The ASN database path is required through `MAXMIND_ASN_DB`.
- You must provide at least one lookup input: `-i` (IP) or `-d` (domain).
- If both `-i` and `-d` are provided, `-i` is used.
- The Makefile build uses `-trimpath -ldflags="-s -w"` to reduce binary size.
- Output is printed in two sections:
	- `---- Geo Lookup ----` as a two-column table.
	- `---- ASN Lookup ----` as `AS<number> <organization>`.
- On invalid input, the tool exits with a non-zero status and prints a validation error.
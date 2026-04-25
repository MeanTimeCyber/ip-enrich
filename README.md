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

Set the MaxMind database paths with `MAXMIND_CITY_DB` and `MAXMIND_ASN_DB`, then pass one of these inputs:

- `-i <ip-address>` for a single IP
- `-d <domain>` for a single domain
- `-il <file>` for a file of IPs (one per line)
- `-dl <file>` for a file of domains (one per line)

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

Look up IPs from a list file:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -il ip_list.txt
```

Look up domains from a list file:

```bash
MAXMIND_CITY_DB=sources/GeoLite2-City.mmdb MAXMIND_ASN_DB=sources/GeoLite2-ASN.mmdb ./ip-enrich -dl domain_list.txt
```

## Usage

```text
./ip-enrich [-i <ip-address> | -d <domain> | -il <file> | -dl <file>] [-json] [-dbinfo]
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
- You must provide at least one lookup input: `-i`, `-d`, `-il`, or `-dl`.
- List files for `-il` and `-dl` are read one line at a time; blank lines are skipped.
- If multiple lookup flags are provided, the first matching mode is used in this order: `-d`, then `-dl`, then `-i`, then `-il`.
- The Makefile build uses `-trimpath -ldflags="-s -w"` to reduce binary size.
- Output is printed in two sections:
	- `---- Geo Lookup ----` as a two-column table.
	- `---- ASN Lookup ----` as `AS<number> <organization>`.
- On invalid input, the tool exits with a non-zero status and prints a validation error.
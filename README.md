# ip-enrich

`ip-enrich` is a small Go CLI for enriching an IP address with location data from a MaxMind GeoLite2 City database.

## Requirements

- Go 1.25.8 or newer
- A MaxMind GeoLite2 City database file

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

Set the MaxMind database path with the `MAXMINDCITYDB` environment variable, then pass the IP address with `-i`.

Run directly with Go:

```bash
MAXMINDCITYDB=sources/GeoLite2-City.mmdb go run cli/*.go -i 109.158.10.179
```

Run the built binary:

```bash
MAXMINDCITYDB=sources/GeoLite2-City.mmdb ./ip-enrich -i 109.158.10.179
```

Or use the Makefile run target:

```bash
MAXMINDCITYDB=sources/GeoLite2-City.mmdb make run IP=109.158.10.179
```

## Usage

```text
./ip-enrich -i <ip-address>
```

Example:

```bash
MAXMINDCITYDB=sources/GeoLite2-City.mmdb ./ip-enrich -i 8.8.8.8
```

## Notes

- The database path is required through `MAXMINDCITYDB`.
- The `-i` flag is required.
- The Makefile build uses `-trimpath -ldflags="-s -w"` to reduce binary size.
- Output is printed as a two-column table.
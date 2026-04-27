# Orthocal Release Checklist

Use this checklist before tagging a release.

## Database

- Verify the generator creates a schema version 4 database.
- Run `orthocal info --db ./orthodox-calendar.db`.
- Confirm `schema_version: 4`.
- Confirm `includes_calendar_events: true`.
- Test update with a local database file.

## Tests

```sh
go test ./...
make test
```

## Build

```sh
make clean
make build VERSION=0.1.0
./bin/orthocal version
```

## Runtime Checks

```sh
./bin/orthocal info --db ./orthodox-calendar.db
./bin/orthocal date 2026-04-12 --db ./orthodox-calendar.db
./bin/orthocal events 2026-04-12 --db ./orthodox-calendar.db
./bin/orthocal search events Pascha --db ./orthodox-calendar.db
./bin/orthocal search saints John --db ./orthodox-calendar.db --limit 5
./bin/orthocal export-web /tmp/orthocal-site --db ./orthodox-calendar.db
```

## Install Check

```sh
make install-all DESTDIR=/tmp/orthocal-dest VERSION=0.1.0
test -x /tmp/orthocal-dest/usr/local/bin/orthocal
test -f /tmp/orthocal-dest/usr/local/man/man1/orthocal.1
```

## Packaging

- Package on OpenBSD.
- Package on FreeBSD.
- Package on Linux.
- Verify shell completions install paths.
- Verify the OpenBSD rc.d script can start the local server with a configured DB path.

## Release

- Tag the release.
- Attach binary artifacts.
- Attach database artifacts.
- Attach checksums.

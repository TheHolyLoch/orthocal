# Orthocal

Orthocal is a native POSIX command line program for reading Orthodox Old Style calendar data from a local SQLite database.

Schema version 4 databases include calculated calendar events, fasting seasons, fast-free periods, remembrances, and calculated fasting levels.

## Table of Contents
- [Requirements](#requirements)
- [Setup](#setup)
- [Building from Source](#building-from-source)
- [Usage](#usage)
- [Database Path](#database-path)
- [Update Behavior](#update-behavior)
- [Schema Compatibility](#schema-compatibility)
- [Shell Completions](#shell-completions)
- [Manual Page](#manual-page)
- [OpenBSD Notes](#openbsd-notes)
- [Examples](#examples)

## Requirements
- [Go](https://go.dev)

## Setup
```sh
go build ./cmd/orthocal
```

## Building from Source
Build with Go directly:
```sh
go build -o bin/orthocal ./cmd/orthocal
```

Or use the Makefile:
```sh
make build
make build VERSION=0.1.0
make test
make install
```

Makefile targets:

| Target                | Description                                |
| --------------------- | ------------------------------------------ |
| `build`               | Build `bin/orthocal`                       |
| `test`                | Run `go test ./...`                        |
| `fmt`                 | Run `gofmt` on Go files                    |
| `clean`               | Remove `bin/`                              |
| `install`             | Install to `$(DESTDIR)$(BINDIR)`           |
| `install-completions` | Install bash, zsh, and fish completions    |
| `install-man`         | Install `docs/orthocal.1`                  |
| `install-all`         | Install binary, completions, and man page  |
| `uninstall`           | Remove installed binary, completions, page |

## Usage
```sh
orthocal [--db PATH] [--plain] [--json] COMMAND [ARGS]
```

| Command                 | Description                          |
| ----------------------- | ------------------------------------ |
| `today`                 | Show the local system date           |
| `tomorrow`              | Show the local system date plus one  |
| `date YYYY-MM-DD`       | Show a specific Gregorian date       |
| `saints YYYY-MM-DD`     | Show saints for a Gregorian date     |
| `readings YYYY-MM-DD`   | Show readings for a Gregorian date   |
| `hymns YYYY-MM-DD`      | Show hymns for a Gregorian date      |
| `search CATEGORY QUERY` | Search calendar data                 |
| `serve`                 | Start the local read-only web server |
| `export-web OUTPUT_DIR` | Export a static read-only website    |
| `info`                  | Show database metadata and counts    |
| `update SOURCE`         | Replace the configured database      |
| `version`               | Show version and build information   |

Show database information:
```sh
orthocal info
orthocal info --db ./orthodox-calendar.db
orthocal version
```

Show calendar days:
```sh
orthocal today --db ./orthodox-calendar.db
orthocal tomorrow --db ./orthodox-calendar.db
orthocal date 2026-04-12 --db ./orthodox-calendar.db
orthocal saints 2026-04-12 --db ./orthodox-calendar.db
orthocal readings 2026-04-12 --db ./orthodox-calendar.db
orthocal hymns 2026-04-12 --db ./orthodox-calendar.db
```

Search calendar data:
```sh
orthocal search saints John --db ./orthodox-calendar.db
orthocal search western Osburga --db ./orthodox-calendar.db
orthocal search primary Climacus --db ./orthodox-calendar.db
orthocal search readings "John 20" --db ./orthodox-calendar.db
orthocal search hymns resurrection --db ./orthodox-calendar.db
orthocal search events Pascha --db ./orthodox-calendar.db
orthocal search feasts Pascha --db ./orthodox-calendar.db
orthocal search fasts Lent --db ./orthodox-calendar.db
orthocal search remembrances departed --db ./orthodox-calendar.db
orthocal search saints John --limit 50 --db ./orthodox-calendar.db
```

Print JSON:
```sh
orthocal date 2026-04-12 --db ./orthodox-calendar.db --json
orthocal saints 2026-04-12 --db ./orthodox-calendar.db --json
orthocal readings 2026-04-12 --db ./orthodox-calendar.db --json
orthocal hymns 2026-04-12 --db ./orthodox-calendar.db --json
orthocal search saints John --db ./orthodox-calendar.db --json
orthocal info --db ./orthodox-calendar.db --json
```

Run the local web server:
```sh
orthocal serve --db ./orthodox-calendar.db
orthocal serve --db ./orthodox-calendar.db --addr 127.0.0.1:9090
```

Server mode is read-only and binds to `127.0.0.1:8080` by default.

Export a static website:
```sh
orthocal export-web ./site --db ./orthodox-calendar.db
```

The export is static and read-only. Exported JSON is written under `api/`.

API endpoints:

| Endpoint                   | Description                  |
| -------------------------- | ---------------------------- |
| `/api/today`               | Today's calendar data        |
| `/api/tomorrow`            | Tomorrow's calendar data     |
| `/api/date/YYYY-MM-DD`     | Calendar data for a date     |
| `/api/saints/YYYY-MM-DD`   | Saints for a date            |
| `/api/readings/YYYY-MM-DD` | Readings for a date          |
| `/api/hymns/YYYY-MM-DD`    | Hymns for a date             |
| `/api/info`                | Database metadata and counts |

## Database Path
If `--db` is omitted, Orthocal checks paths in this order:

| Order | Path                                           |
| ----- | ---------------------------------------------- |
| 1     | `$ORTHOCAL_DB`                                 |
| 2     | `$XDG_DATA_HOME/orthocal/orthodox-calendar.db` |
| 3     | `~/.local/share/orthocal/orthodox-calendar.db` |

## Update Behavior
`orthocal update SOURCE` accepts a local file path or an `http`/`https` URL. It downloads or copies the source to a temporary file, validates the SQLite database, then atomically replaces the configured database.

Validation requires `PRAGMA integrity_check` to return `ok`, an `app_metadata` table, and an `app_metadata` row with key `schema_version`.

Newer schema versions are rejected unless `--force` is passed:
```sh
orthocal update ./orthodox-calendar.db --force
```

One backup is kept at `<database>.bak`.

## Schema Compatibility
Orthocal supports schema version 4. Older databases still run with degraded event support.

`orthocal info` prints the detected schema version and compatibility message. Newer databases show a warning, but read-only commands continue if the required columns are still compatible.

## Shell Completions
Static completions are included for bash, zsh, and fish.

Install them with:
```sh
make install-completions
```

The zsh completion is installed as `_orthocal`.

## Manual Page
Install the manual page with:
```sh
make install-man
```

Then read it with:
```sh
man orthocal
```

## OpenBSD Notes
Orthocal includes OpenBSD `pledge(2)` and `unveil(2)` support. On non-OpenBSD systems these calls are no-ops.

Hardening notes:
- Read-only CLI commands unveil the configured database read-only and pledge `stdio rpath`.
- Server mode unveils the configured database read-only and pledges `stdio rpath inet dns`.
- Update mode uses a conservative pledge for local and HTTP updates because it needs read, write, create, rename, and metadata operations.
- HTTP update mode does not unveil paths before download so HTTPS certificate loading keeps working.

The web server binds to `127.0.0.1:8080` by default.

OpenBSD rc.d example:
```sh
cp packaging/openbsd/orthocal.rc /etc/rc.d/orthocal
rcctl enable orthocal
rcctl set orthocal flags "serve --db /var/db/orthocal/orthodox-calendar.db --addr 127.0.0.1:8080"
rcctl start orthocal
```

## Examples
```sh
orthocal today
orthocal tomorrow
orthocal date 2026-04-12
orthocal saints 2026-04-12
orthocal readings 2026-04-12
orthocal hymns 2026-04-12
orthocal search saints John
orthocal search western Osburga
orthocal search primary Climacus
orthocal search readings "John 20"
orthocal search hymns resurrection
orthocal search events Pascha
orthocal search feasts Pascha
orthocal search fasts Lent
orthocal search remembrances departed
orthocal search saints John --limit 50
orthocal serve --db ./orthodox-calendar.db
orthocal serve --db ./orthodox-calendar.db --addr 127.0.0.1:9090
orthocal export-web ./site --db ./orthodox-calendar.db
orthocal version
orthocal date 2026-04-12 --db ./orthodox-calendar.db --json
orthocal saints 2026-04-12 --db ./orthodox-calendar.db --json
orthocal search saints John --json
orthocal info --db ./orthodox-calendar.db --json
orthocal --db ./orthodox-calendar.db info
orthocal update ./orthodox-calendar.db
orthocal update ./orthodox-calendar.db --force
orthocal update ./orthodox-calendar.db --db ~/.local/share/orthocal/orthodox-calendar.db
orthocal update https://example.org/orthodox-calendar.db
```

# Orthocal

Orthocal is a native POSIX command line program for reading Orthodox Old Style calendar data from a local SQLite database.

This pass supports database path resolution, SQLite opening, `info`, `today`, `tomorrow`, and `date`.

## Table of Contents
- [Requirements](#requirements)
- [Setup](#setup)
- [Usage](#usage)
- [Database Path](#database-path)
- [Update Behavior](#update-behavior)
- [Examples](#examples)

## Requirements
- [Go](https://go.dev)

## Setup
```sh
go build ./cmd/orthocal
```

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

Show database information:
```sh
orthocal info
orthocal info --db ./orthodox-calendar.db
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

One backup is kept at `<database>.bak`.

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
orthocal search saints John --limit 50
orthocal serve --db ./orthodox-calendar.db
orthocal serve --db ./orthodox-calendar.db --addr 127.0.0.1:9090
orthocal export-web ./site --db ./orthodox-calendar.db
orthocal date 2026-04-12 --db ./orthodox-calendar.db --json
orthocal saints 2026-04-12 --db ./orthodox-calendar.db --json
orthocal search saints John --json
orthocal info --db ./orthodox-calendar.db --json
orthocal --db ./orthodox-calendar.db info
orthocal update ./orthodox-calendar.db
orthocal update ./orthodox-calendar.db --db ~/.local/share/orthocal/orthodox-calendar.db
orthocal update https://example.org/orthodox-calendar.db
```

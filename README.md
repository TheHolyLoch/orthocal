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

| Command              | Description                         |
| -------------------- | ----------------------------------- |
| `today`              | Show the local system date          |
| `tomorrow`           | Show the local system date plus one |
| `date YYYY-MM-DD`    | Show a specific Gregorian date      |
| `info`               | Show database metadata and counts   |
| `update SOURCE`      | Replace the configured database     |

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
```

Print JSON:
```sh
orthocal date 2026-04-12 --db ./orthodox-calendar.db --json
orthocal info --db ./orthodox-calendar.db --json
```

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
orthocal date 2026-04-12 --db ./orthodox-calendar.db --json
orthocal info --db ./orthodox-calendar.db --json
orthocal --db ./orthodox-calendar.db info
orthocal update ./orthodox-calendar.db
orthocal update ./orthodox-calendar.db --db ~/.local/share/orthocal/orthodox-calendar.db
orthocal update https://example.org/orthodox-calendar.db
```

# Orthocal

Orthocal is a native POSIX command line program for reading Orthodox Old Style calendar data from a local SQLite database.

This pass creates the CLI skeleton only. Database reads and update verification are added in later passes.

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

## Database Path
`--db PATH` is parsed now. Default path resolution is added in the next pass.

## Update Behavior
`orthocal update SOURCE` validates that `SOURCE` was provided and prints a placeholder in this pass.

## Examples
```sh
orthocal today
orthocal tomorrow
orthocal date 2026-04-12
orthocal --db ./orthodox-calendar.db info
orthocal update ./orthodox-calendar.db
```

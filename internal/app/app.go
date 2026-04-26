// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/app/app.go

package app

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"orthocal/internal/config"
	"orthocal/internal/db"
)

const usageText = `orthocal [--db PATH] [--plain] [--json] COMMAND [ARGS]

Commands:
  today              Show the local system date
  tomorrow           Show the local system date plus one
  date YYYY-MM-DD    Show a specific Gregorian date
  info               Show database metadata and counts
  update SOURCE      Replace the configured database

Options:
  --db PATH          Use a specific SQLite database
  --plain            Disable terminal styling
  --json             Print JSON for today, tomorrow, and date
  --help             Print help text
`

type options struct {
	dbPath string
	help   bool
	json   bool
	plain  bool
}

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	opts, rest, err := parse_options(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	if opts.help {
		print_usage(stdout)
		return 0
	}

	if len(rest) == 0 {
		print_usage(stderr)
		return 2
	}

	command := rest[0]
	commandArgs := rest[1:]

	switch command {
	case "today":
		return run_today(commandArgs, stdout, stderr)

	case "tomorrow":
		return run_tomorrow(commandArgs, stdout, stderr)

	case "date":
		return run_date(commandArgs, stdout, stderr)

	case "info":
		return run_info(opts, commandArgs, stdout, stderr)

	case "update":
		return run_update(commandArgs, stdout, stderr)

	case "help":
		print_usage(stdout)
		return 0
	}

	fmt.Fprintf(stderr, "error: unknown command %q\n", command)
	return 2
}

func parse_options(args []string) (options, []string, error) {
	opts := options{}
	rest := []string{}

	for index := 0; index < len(args); index++ {
		arg := args[index]

		switch {
		case arg == "--db":
			if index+1 >= len(args) {
				return options{}, nil, errors.New("error: --db requires PATH")
			}

			index++
			opts.dbPath = args[index]

		case strings.HasPrefix(arg, "--db="):
			opts.dbPath = strings.TrimPrefix(arg, "--db=")
			if opts.dbPath == "" {
				return options{}, nil, errors.New("error: --db requires PATH")
			}

		case arg == "--plain":
			opts.plain = true

		case arg == "--json":
			opts.json = true

		case arg == "--help":
			opts.help = true

		default:
			rest = append(rest, arg)
		}
	}

	return opts, rest, nil
}

func print_usage(output io.Writer) {
	fmt.Fprint(output, usageText)
}

func run_date(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "error: date requires YYYY-MM-DD")
		return 2
	}

	if _, err := time.Parse("2006-01-02", args[0]); err != nil {
		fmt.Fprintln(stderr, "error: date must use YYYY-MM-DD")
		return 2
	}

	fmt.Fprintf(stdout, "date lookup for %s will be added in a later pass\n", args[0])
	return 0
}

func run_info(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: info does not accept arguments")
		return 2
	}

	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer conn.Close()

	view := db.InfoView{DatabasePath: dbPath}

	metadata, err := db.MetadataRows(conn)
	if err != nil {
		if errors.Is(err, db.ErrTableMissing) {
			view.MetadataUnavailable = true
		} else {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else {
		view.Metadata = metadata
	}

	if count, err := db.CountRows(conn, "calendar_days"); err == nil {
		view.CalendarDaysCount = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if count, err := db.CountRows(conn, "saints"); err == nil {
		view.SaintsCount = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if count, err := db.CountRows(conn, "scripture_readings"); err == nil {
		view.ScriptureReadingsCount = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if count, err := db.CountRows(conn, "hymns"); err == nil {
		view.HymnsCount = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	print_info(stdout, view)
	return 0
}

func print_info(stdout io.Writer, view db.InfoView) {
	fmt.Fprintf(stdout, "database path: %s\n", view.DatabasePath)

	if view.MetadataUnavailable {
		fmt.Fprintln(stdout, "app_metadata: unavailable")
	} else {
		fmt.Fprintln(stdout, "app_metadata:")
		for _, item := range view.Metadata {
			fmt.Fprintf(stdout, "\t%s: %s\n", item.Key, item.Value)
		}
	}

	fmt.Fprintf(stdout, "calendar_days: %d\n", view.CalendarDaysCount)
	fmt.Fprintf(stdout, "saints: %d\n", view.SaintsCount)
	fmt.Fprintf(stdout, "scripture_readings: %d\n", view.ScriptureReadingsCount)
	fmt.Fprintf(stdout, "hymns: %d\n", view.HymnsCount)
}

func run_today(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: today does not accept arguments")
		return 2
	}

	fmt.Fprintln(stdout, "date lookup will be added in a later pass")
	return 0
}

func run_tomorrow(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: tomorrow does not accept arguments")
		return 2
	}

	fmt.Fprintln(stdout, "date lookup will be added in a later pass")
	return 0
}

func run_update(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "error: update requires SOURCE")
		return 2
	}

	fmt.Fprintf(stdout, "update support for %s will be added in a later pass\n", args[0])
	return 0
}

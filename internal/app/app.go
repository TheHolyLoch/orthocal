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
	"orthocal/internal/render"
	"orthocal/internal/update"
)

const usageText = `orthocal [--db PATH] [--plain] [--json] COMMAND [ARGS]

Commands:
  today              Show the local system date
  tomorrow           Show the local system date plus one
  date YYYY-MM-DD    Show a specific Gregorian date
  saints YYYY-MM-DD  Show saints for a Gregorian date
  readings YYYY-MM-DD
                     Show scripture readings for a Gregorian date
  hymns YYYY-MM-DD   Show hymns for a Gregorian date
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
		return run_today(opts, commandArgs, stdout, stderr)

	case "tomorrow":
		return run_tomorrow(opts, commandArgs, stdout, stderr)

	case "date":
		return run_date(opts, commandArgs, stdout, stderr)

	case "hymns":
		return run_hymns(opts, commandArgs, stdout, stderr)

	case "info":
		return run_info(opts, commandArgs, stdout, stderr)

	case "readings":
		return run_readings(opts, commandArgs, stdout, stderr)

	case "saints":
		return run_saints(opts, commandArgs, stdout, stderr)

	case "update":
		return run_update(opts, commandArgs, stdout, stderr)

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

func open_day_view(opts options, value string) (db.DayView, bool, error) {
	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		return db.DayView{}, false, err
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		return db.DayView{}, false, err
	}
	defer conn.Close()

	return db.DayViewByGregorianDate(conn, value)
}

func open_hymns_view(opts options, value string) (db.HymnsView, bool, error) {
	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		return db.HymnsView{}, false, err
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		return db.HymnsView{}, false, err
	}
	defer conn.Close()

	return db.HymnsViewByGregorianDate(conn, value)
}

func open_readings_view(opts options, value string) (db.ReadingsView, bool, error) {
	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		return db.ReadingsView{}, false, err
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		return db.ReadingsView{}, false, err
	}
	defer conn.Close()

	return db.ReadingsViewByGregorianDate(conn, value)
}

func open_saints_view(opts options, value string) (db.SaintsView, bool, error) {
	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		return db.SaintsView{}, false, err
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		return db.SaintsView{}, false, err
	}
	defer conn.Close()

	return db.SaintsViewByGregorianDate(conn, value)
}

func render_day(opts options, value string, stdout io.Writer, stderr io.Writer) int {
	view, found, err := open_day_view(opts, value)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if !found {
		fmt.Fprintf(stderr, "date not found: %s\n", value)
		return 1
	}

	if opts.json {
		if err := render.RenderDayJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.Day(stdout, view)
	return 0
}

func run_date(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	value, ok := parse_date_argument(args, "date", stderr)
	if !ok {
		return 2
	}

	return render_day(opts, value, stdout, stderr)
}

func run_hymns(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	value, ok := parse_date_argument(args, "hymns", stderr)
	if !ok {
		return 2
	}

	view, found, err := open_hymns_view(opts, value)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if !found {
		fmt.Fprintf(stderr, "date not found: %s\n", value)
		return 1
	}

	if opts.json {
		if err := render.RenderHymnsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.Hymns(stdout, view)
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

	view := db.InfoView{
		DatabasePath: dbPath,
		Metadata:     []db.Metadata{},
	}

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
		view.Counts.CalendarDays = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if count, err := db.CountRows(conn, "saints"); err == nil {
		view.Counts.Saints = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if count, err := db.CountRows(conn, "scripture_readings"); err == nil {
		view.Counts.ScriptureReadings = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if count, err := db.CountRows(conn, "hymns"); err == nil {
		view.Counts.Hymns = count
	} else {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if opts.json {
		if err := render.RenderInfoJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
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

	fmt.Fprintf(stdout, "calendar_days: %d\n", view.Counts.CalendarDays)
	fmt.Fprintf(stdout, "saints: %d\n", view.Counts.Saints)
	fmt.Fprintf(stdout, "scripture_readings: %d\n", view.Counts.ScriptureReadings)
	fmt.Fprintf(stdout, "hymns: %d\n", view.Counts.Hymns)
}

func run_readings(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	value, ok := parse_date_argument(args, "readings", stderr)
	if !ok {
		return 2
	}

	view, found, err := open_readings_view(opts, value)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if !found {
		fmt.Fprintf(stderr, "date not found: %s\n", value)
		return 1
	}

	if opts.json {
		if err := render.RenderReadingsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.Readings(stdout, view)
	return 0
}

func run_saints(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	value, ok := parse_date_argument(args, "saints", stderr)
	if !ok {
		return 2
	}

	view, found, err := open_saints_view(opts, value)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if !found {
		fmt.Fprintf(stderr, "date not found: %s\n", value)
		return 1
	}

	if opts.json {
		if err := render.RenderSaintsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.Saints(stdout, view)
	return 0
}

func run_today(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: today does not accept arguments")
		return 2
	}

	value := time.Now().Format("2006-01-02")
	return render_day(opts, value, stdout, stderr)
}

func run_tomorrow(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: tomorrow does not accept arguments")
		return 2
	}

	value := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	return render_day(opts, value, stdout, stderr)
}

func run_update(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "error: update requires SOURCE")
		return 2
	}

	targetPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	result, err := update.UpdateDatabase(targetPath, args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintf(stdout, "source: %s\n", result.Source)
	fmt.Fprintf(stdout, "target path: %s\n", result.TargetPath)
	if result.BackupCreated {
		fmt.Fprintf(stdout, "backup path: %s\n", result.BackupPath)
	}
	fmt.Fprintln(stdout, "validation: ok")
	return 0
}

func parse_date_argument(args []string, command string, stderr io.Writer) (string, bool) {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "error: %s requires YYYY-MM-DD\n", command)
		return "", false
	}

	if _, err := time.Parse("2006-01-02", args[0]); err != nil {
		fmt.Fprintln(stderr, "error: date must use YYYY-MM-DD")
		return "", false
	}

	return args[0], true
}

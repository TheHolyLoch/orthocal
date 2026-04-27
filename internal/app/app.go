// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/app/app.go

package app

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"orthocal/internal/config"
	"orthocal/internal/db"
	"orthocal/internal/platform"
	"orthocal/internal/render"
	"orthocal/internal/server"
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
  search CATEGORY QUERY
                     Search saints, western, primary, readings, hymns, or events
  serve              Start the local read-only web server
  export-web OUTPUT_DIR
                     Export a static read-only website
  info               Show database metadata and counts
  update SOURCE      Replace the configured database

Options:
  --db PATH          Use a specific SQLite database
  --addr ADDRESS     Web server address, default 127.0.0.1:8080
  --plain            Disable terminal styling
  --json             Print JSON for supported commands
  --limit N          Limit search results, default 25, max 200
  --help             Print help text
`

type options struct {
	addr   string
	dbPath string
	help   bool
	json   bool
	limit  int
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

	case "export-web":
		return run_export_web(opts, commandArgs, stdout, stderr)

	case "info":
		return run_info(opts, commandArgs, stdout, stderr)

	case "readings":
		return run_readings(opts, commandArgs, stdout, stderr)

	case "saints":
		return run_saints(opts, commandArgs, stdout, stderr)

	case "search":
		return run_search(opts, commandArgs, stdout, stderr)

	case "serve":
		return run_serve(opts, commandArgs, stdout, stderr)

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
	opts := options{
		addr:  "127.0.0.1:8080",
		limit: 25,
	}
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

		case arg == "--addr":
			if index+1 >= len(args) {
				return options{}, nil, errors.New("error: --addr requires ADDRESS")
			}

			index++
			opts.addr = args[index]

		case strings.HasPrefix(arg, "--addr="):
			opts.addr = strings.TrimPrefix(arg, "--addr=")
			if opts.addr == "" {
				return options{}, nil, errors.New("error: --addr requires ADDRESS")
			}

		case arg == "--plain":
			opts.plain = true

		case arg == "--json":
			opts.json = true

		case arg == "--limit":
			if index+1 >= len(args) {
				return options{}, nil, errors.New("error: --limit requires N")
			}

			index++
			limit, err := strconv.Atoi(args[index])
			if err != nil {
				return options{}, nil, errors.New("error: --limit must be a number")
			}

			opts.limit = db.ClampLimit(limit)

		case strings.HasPrefix(arg, "--limit="):
			value := strings.TrimPrefix(arg, "--limit=")
			limit, err := strconv.Atoi(value)
			if err != nil {
				return options{}, nil, errors.New("error: --limit must be a number")
			}

			opts.limit = db.ClampLimit(limit)

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

func harden_export_web(dbPath string, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	if err := platform.UnveilReadOnly(dbPath); err != nil {
		return err
	}

	if err := platform.UnveilReadWrite(outputDir); err != nil {
		return err
	}

	if err := platform.UnveilNone(); err != nil {
		return err
	}

	return platform.PledgeNetworkClient()
}

func harden_read_only_cli(dbPath string) error {
	if err := platform.UnveilReadOnly(dbPath); err != nil {
		return err
	}

	if err := platform.UnveilNone(); err != nil {
		return err
	}

	return platform.PledgeCLI()
}

func harden_server(dbPath string) error {
	if err := platform.UnveilReadOnly(dbPath); err != nil {
		return err
	}

	if err := platform.UnveilNone(); err != nil {
		return err
	}

	return platform.PledgeServer()
}

func harden_update(targetPath string, source string) error {
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	if update.IsHTTPSource(source) {
		return platform.PledgeNetworkClient()
	}

	if err := platform.UnveilReadOnly(source); err != nil {
		return err
	}

	if err := platform.UnveilReadWrite(targetDir); err != nil {
		return err
	}

	if err := platform.UnveilNone(); err != nil {
		return err
	}

	return platform.PledgeNetworkClient()
}

func open_day_view(opts options, value string) (db.DayView, bool, error) {
	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		return db.DayView{}, false, err
	}

	if err := harden_read_only_cli(dbPath); err != nil {
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

	if err := harden_read_only_cli(dbPath); err != nil {
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

	if err := harden_read_only_cli(dbPath); err != nil {
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

	if err := harden_read_only_cli(dbPath); err != nil {
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

func run_export_web(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "error: export-web requires OUTPUT_DIR")
		return 2
	}

	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if err := harden_export_web(dbPath, args[0]); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer conn.Close()

	webServer, err := server.New(conn, server.Config{
		DatabasePath: dbPath,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	count, err := webServer.ExportWeb(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintf(stdout, "output directory: %s\n", args[0])
	fmt.Fprintf(stdout, "days exported: %d\n", count)
	return 0
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

	if err := harden_read_only_cli(dbPath); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer conn.Close()

	view, err := db.InfoViewByPath(conn, dbPath)
	if err != nil {
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
	fmt.Fprintf(stdout, "calendar_events: %d\n", view.Counts.CalendarEvents)
	fmt.Fprintf(stdout, "calendar_day_events: %d\n", view.Counts.CalendarDayEvents)
	fmt.Fprintf(stdout, "saints: %d\n", view.Counts.Saints)
	fmt.Fprintf(stdout, "scripture_readings: %d\n", view.Counts.ScriptureReadings)
	fmt.Fprintf(stdout, "hymns: %d\n", view.Counts.Hymns)
	if view.SchemaNote != "" {
		fmt.Fprintf(stdout, "note: %s\n", view.SchemaNote)
	}
}

func run_search(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "error: search requires CATEGORY and QUERY")
		return 2
	}

	category := args[0]
	query := strings.TrimSpace(strings.Join(args[1:], " "))
	if query == "" {
		fmt.Fprintln(stderr, "error: search query cannot be empty")
		return 2
	}

	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if err := harden_read_only_cli(dbPath); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer conn.Close()

	switch category {
	case "saints":
		return run_search_saints(conn, opts, query, "saints", false, false, stdout, stderr)

	case "western":
		return run_search_saints(conn, opts, query, "western", true, false, stdout, stderr)

	case "primary":
		return run_search_saints(conn, opts, query, "primary", false, true, stdout, stderr)

	case "readings":
		return run_search_readings(conn, opts, query, stdout, stderr)

	case "hymns":
		return run_search_hymns(conn, opts, query, stdout, stderr)

	case "events", "feasts", "fasts", "remembrances":
		return run_search_events(conn, opts, query, category, stdout, stderr)
	}

	fmt.Fprintf(stderr, "error: unknown search category %q\n", category)
	return 2
}

func run_search_events(conn *sql.DB, opts options, query string, category string, stdout io.Writer, stderr io.Writer) int {
	results, err := db.SearchEvents(conn, query, category, opts.limit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	view := db.SearchEventsView{
		Query:    query,
		Category: category,
		Limit:    db.ClampLimit(opts.limit),
		Results:  results,
	}

	if opts.json {
		if err := render.RenderSearchEventsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.SearchEvents(stdout, view)
	return 0
}

func run_search_hymns(conn *sql.DB, opts options, query string, stdout io.Writer, stderr io.Writer) int {
	results, err := db.SearchHymns(conn, query, opts.limit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	view := db.SearchHymnsView{
		Query:    query,
		Category: "hymns",
		Limit:    db.ClampLimit(opts.limit),
		Results:  results,
	}

	if opts.json {
		if err := render.RenderSearchHymnsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.SearchHymns(stdout, view)
	return 0
}

func run_search_readings(conn *sql.DB, opts options, query string, stdout io.Writer, stderr io.Writer) int {
	results, err := db.SearchReadings(conn, query, opts.limit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	view := db.SearchReadingsView{
		Query:    query,
		Category: "readings",
		Limit:    db.ClampLimit(opts.limit),
		Results:  results,
	}

	if opts.json {
		if err := render.RenderSearchReadingsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.SearchReadings(stdout, view)
	return 0
}

func run_search_saints(conn *sql.DB, opts options, query string, category string, westernOnly bool, primaryOnly bool, stdout io.Writer, stderr io.Writer) int {
	results, err := db.SearchSaints(conn, query, westernOnly, primaryOnly, opts.limit)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	view := db.SearchSaintsView{
		Query:    query,
		Category: category,
		Limit:    db.ClampLimit(opts.limit),
		Results:  results,
	}

	if opts.json {
		if err := render.RenderSearchSaintsJSON(view); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}

		return 0
	}

	render.SearchSaints(stdout, view)
	return 0
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

func run_serve(opts options, args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: serve does not accept arguments")
		return 2
	}

	dbPath, err := config.ResolveDBPath(opts.dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if err := harden_server(dbPath); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer conn.Close()

	if _, err := db.InfoViewByPath(conn, dbPath); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	webServer, err := server.New(conn, server.Config{
		Addr:         opts.addr,
		DatabasePath: dbPath,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmt.Fprintf(stdout, "serving address: %s\n", opts.addr)
	fmt.Fprintf(stdout, "database path: %s\n", dbPath)

	if err := webServer.Serve(opts.addr); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

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

	if err := harden_update(targetPath, args[0]); err != nil {
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

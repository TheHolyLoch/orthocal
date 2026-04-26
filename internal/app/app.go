// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/app/app.go

package app

import (
	"flag"
	"fmt"
	"io"
	"time"
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
	opts, rest, err := parse_options(args, stderr)
	if err != nil {
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
		return run_info(commandArgs, stdout, stderr)

	case "update":
		return run_update(commandArgs, stdout, stderr)

	case "help":
		print_usage(stdout)
		return 0
	}

	fmt.Fprintf(stderr, "error: unknown command %q\n", command)
	return 2
}

func parse_options(args []string, output io.Writer) (options, []string, error) {
	opts := options{}

	flags := flag.NewFlagSet("orthocal", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.dbPath, "db", "", "Use a specific SQLite database")
	flags.BoolVar(&opts.plain, "plain", false, "Disable terminal styling")
	flags.BoolVar(&opts.json, "json", false, "Print JSON for today, tomorrow, and date")
	flags.BoolVar(&opts.help, "help", false, "Print help text")
	flags.Usage = func() {
		print_usage(output)
	}

	if err := flags.Parse(args); err != nil {
		return options{}, nil, err
	}

	return opts, flags.Args(), nil
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

func run_info(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "error: info does not accept arguments")
		return 2
	}

	fmt.Fprintln(stdout, "DB support will be added in the next pass")
	return 0
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

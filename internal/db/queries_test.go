// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/queries_test.go

package db

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestClampLimit(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "default", in: 0, want: 25},
		{name: "negative", in: -4, want: 25},
		{name: "inside", in: 50, want: 50},
		{name: "max", in: 250, want: 200},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ClampLimit(test.in); got != test.want {
				t.Fatalf("expected %d, got %d", test.want, got)
			}
		})
	}
}

func TestDayQueries(t *testing.T) {
	conn := test_db(t)
	defer conn.Close()

	day, found, err := DayByGregorianDate(conn, "2026-04-12")
	if err != nil {
		t.Fatalf("DayByGregorianDate returned error: %v", err)
	}
	if !found {
		t.Fatal("expected day to be found")
	}
	if day.ID != 1 || day.JulianDate != "2026-03-30" {
		t.Fatalf("unexpected day: %#v", day)
	}

	missing, found, err := DayByGregorianDate(conn, "2026-04-13")
	if err != nil {
		t.Fatalf("DayByGregorianDate missing returned error: %v", err)
	}
	if found {
		t.Fatalf("expected missing date, got %#v", missing)
	}

	view, found, err := DayViewByGregorianDate(conn, "2026-04-12")
	if err != nil {
		t.Fatalf("DayViewByGregorianDate returned error: %v", err)
	}
	if !found {
		t.Fatal("expected day view to be found")
	}
	if len(view.PrimarySaints) != 2 || len(view.WesternSaints) != 1 || len(view.ScriptureReadings) != 3 || len(view.Saints) != 3 {
		t.Fatalf("unexpected day view counts: %#v", view)
	}
}

func TestEscapeLikeQuery(t *testing.T) {
	got := EscapeLikeQuery(`100%\_`)
	want := `100\%\\\_`
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFocusedViews(t *testing.T) {
	conn := test_db(t)
	defer conn.Close()

	saints, found, err := SaintsViewByGregorianDate(conn, "2026-04-12")
	if err != nil {
		t.Fatalf("SaintsViewByGregorianDate returned error: %v", err)
	}
	if !found || len(saints.Saints) != 3 {
		t.Fatalf("unexpected saints view: found=%v view=%#v", found, saints)
	}

	readings, found, err := ReadingsViewByGregorianDate(conn, "2026-04-12")
	if err != nil {
		t.Fatalf("ReadingsViewByGregorianDate returned error: %v", err)
	}
	if !found || len(readings.ScriptureReadings) != 3 {
		t.Fatalf("unexpected readings view: found=%v view=%#v", found, readings)
	}

	hymns, found, err := HymnsViewByGregorianDate(conn, "2026-04-12")
	if err != nil {
		t.Fatalf("HymnsViewByGregorianDate returned error: %v", err)
	}
	if !found || len(hymns.Hymns) != 3 {
		t.Fatalf("unexpected hymns view: found=%v view=%#v", found, hymns)
	}
}

func TestListQueries(t *testing.T) {
	conn := test_db(t)
	defer conn.Close()

	saints, err := SaintsByDayID(conn, 1)
	if err != nil {
		t.Fatalf("SaintsByDayID returned error: %v", err)
	}
	if len(saints) != 3 || !saints[0].IsPrimary || !saints[2].IsWestern {
		t.Fatalf("unexpected saints: %#v", saints)
	}

	readings, err := ScriptureByDayID(conn, 1)
	if err != nil {
		t.Fatalf("ScriptureByDayID returned error: %v", err)
	}
	if len(readings) != 3 || readings[2].Description != "Vespers, Gospel" {
		t.Fatalf("unexpected readings: %#v", readings)
	}

	hymns, err := HymnsByDayID(conn, 1)
	if err != nil {
		t.Fatalf("HymnsByDayID returned error: %v", err)
	}
	if len(hymns) != 3 || hymns[0].Title != "Pascha Troparion" {
		t.Fatalf("unexpected hymns: %#v", hymns)
	}
}

func TestSearchQueries(t *testing.T) {
	conn := test_db(t)
	defer conn.Close()

	saints, err := SearchSaints(conn, "John", false, false, 25)
	if err != nil {
		t.Fatalf("SearchSaints returned error: %v", err)
	}
	if len(saints) != 1 {
		t.Fatalf("expected 1 saint result, got %#v", saints)
	}

	western, err := SearchSaints(conn, "Osburga", true, false, 25)
	if err != nil {
		t.Fatalf("SearchSaints western returned error: %v", err)
	}
	if len(western) != 1 || !western[0].IsWestern {
		t.Fatalf("unexpected western results: %#v", western)
	}

	primary, err := SearchSaints(conn, "John", false, true, 25)
	if err != nil {
		t.Fatalf("SearchSaints primary returned error: %v", err)
	}
	if len(primary) != 1 || !primary[0].IsPrimary {
		t.Fatalf("unexpected primary results: %#v", primary)
	}

	readings, err := SearchReadings(conn, "John 20", 25)
	if err != nil {
		t.Fatalf("SearchReadings returned error: %v", err)
	}
	if len(readings) != 1 || readings[0].Description != "Vespers, Gospel" {
		t.Fatalf("unexpected reading results: %#v", readings)
	}

	hymns, err := SearchHymns(conn, "resurrection", 25)
	if err != nil {
		t.Fatalf("SearchHymns returned error: %v", err)
	}
	if len(hymns) != 1 || hymns[0].Title != "Pascha Troparion" {
		t.Fatalf("unexpected hymn results: %#v", hymns)
	}

	literal, err := SearchHymns(conn, "100%", 25)
	if err != nil {
		t.Fatalf("SearchHymns literal returned error: %v", err)
	}
	if len(literal) != 1 || literal[0].Title != "Literal Percent Hymn" {
		t.Fatalf("unexpected literal results: %#v", literal)
	}
}

func TestSchemaVersion(t *testing.T) {
	conn := schema_test_db(t, "4", true)
	defer conn.Close()

	version, known, err := SchemaVersion(conn)
	if err != nil {
		t.Fatalf("SchemaVersion returned error: %v", err)
	}
	if !known || version != 4 {
		t.Fatalf("unexpected schema version: version=%d known=%v", version, known)
	}
}

func TestCompatibilityStatus(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		metadata    bool
		known       bool
		older       bool
		newer       bool
		version     int
		messagePart string
	}{
		{name: "missing", metadata: false, known: false, messagePart: "unknown"},
		{name: "older", value: "3", metadata: true, known: true, older: true, version: 3, messagePart: "older"},
		{name: "current", value: "4", metadata: true, known: true, version: 4, messagePart: "compatible"},
		{name: "newer", value: "5", metadata: true, known: true, newer: true, version: 5, messagePart: "newer"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn := schema_test_db(t, test.value, test.metadata)
			defer conn.Close()

			status, err := SchemaCompatibilityStatus(conn)
			if err != nil {
				t.Fatalf("SchemaCompatibilityStatus returned error: %v", err)
			}

			if status.SchemaKnown != test.known || status.IsOlder != test.older || status.IsNewer != test.newer || status.SchemaVersion != test.version {
				t.Fatalf("unexpected status: %#v", status)
			}

			if !strings.Contains(status.Message, test.messagePart) {
				t.Fatalf("expected message to contain %q, got %q", test.messagePart, status.Message)
			}
		})
	}
}

func test_db(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "orthocal.db")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}

	statements := []string{
		`CREATE TABLE app_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL)`,
		`INSERT INTO app_metadata (key, value) VALUES ('schema_version', '4')`,
		`CREATE TABLE calendar_days (
			id INTEGER PRIMARY KEY,
			dataheader TEXT NOT NULL,
			gregorian_date TEXT NOT NULL UNIQUE,
			gregorian_weekday TEXT NOT NULL,
			julian_date TEXT NOT NULL,
			headerheader TEXT NOT NULL,
			fasting_rule TEXT NOT NULL
		)`,
		`CREATE TABLE saints (
			day_id INTEGER NOT NULL,
			saint_order INTEGER NOT NULL,
			name TEXT NOT NULL,
			icon_file TEXT NOT NULL,
			is_primary INTEGER NOT NULL,
			is_western INTEGER NOT NULL,
			service_rank_code TEXT NOT NULL,
			service_rank_name TEXT NOT NULL
		)`,
		`CREATE TABLE scripture_readings (
			day_id INTEGER NOT NULL,
			reading_order INTEGER NOT NULL,
			verse_reference TEXT NOT NULL,
			description TEXT NOT NULL,
			reading_url TEXT NOT NULL,
			display_text TEXT NOT NULL
		)`,
		`CREATE TABLE hymns (
			day_id INTEGER NOT NULL,
			hymn_order INTEGER NOT NULL,
			section_order INTEGER NOT NULL,
			hymn_type TEXT NOT NULL,
			tone TEXT NOT NULL,
			title TEXT NOT NULL,
			text TEXT NOT NULL
		)`,
		`INSERT INTO calendar_days (
			id,
			dataheader,
			gregorian_date,
			gregorian_weekday,
			julian_date,
			headerheader,
			fasting_rule
		) VALUES (
			1,
			'Sunday April 12, 2026 / March 30, 2026',
			'2026-04-12',
			'Sunday',
			'2026-03-30',
			'The Bright Resurrection of Christ, The Pascha of the Lord.',
			'The End of the Great Lent.'
		)`,
		`INSERT INTO saints VALUES (1, 1, 'Venerable John Climacus of Sinai.', '1.gif', 1, 0, '1', 'without a sign')`,
		`INSERT INTO saints VALUES (1, 2, 'St. Sophronius, bishop of Irkutsk.', '1.gif', 1, 0, '1', 'without a sign')`,
		`INSERT INTO saints VALUES (1, 3, 'St. Osburga of Coventry.', 'o.gif', 0, 1, 'o', 'ordinary/minor')`,
		`INSERT INTO scripture_readings VALUES (1, 1, 'Acts 1:1-8', '', '', 'Acts 1:1-8')`,
		`INSERT INTO scripture_readings VALUES (1, 2, 'John 1:1-17', '', '', 'John 1:1-17')`,
		`INSERT INTO scripture_readings VALUES (1, 3, 'John 20:19-25', 'Vespers, Gospel', '', 'John 20:19-25 - Vespers, Gospel')`,
		`INSERT INTO hymns VALUES (1, 1, 1, 'Troparion', 'V', 'Pascha Troparion', 'Christ is risen from the dead, trampling down death by death. The resurrection is proclaimed.')`,
		`INSERT INTO hymns VALUES (1, 2, 1, 'Kontakion', 'VIII', 'Kontakion', 'Though Thou didst descend into the grave.')`,
		`INSERT INTO hymns VALUES (1, 3, 2, 'Troparion', '', 'Literal Percent Hymn', 'This line contains 100% literal text and under_score text.')`,
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			conn.Close()
			t.Fatalf("failed statement %q: %v", statement, err)
		}
	}

	return conn
}

func schema_test_db(t *testing.T, schemaVersion string, metadata bool) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "schema.db")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}

	if metadata {
		if _, err := conn.Exec(`CREATE TABLE app_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL)`); err != nil {
			conn.Close()
			t.Fatalf("failed creating app_metadata: %v", err)
		}

		if schemaVersion != "" {
			if _, err := conn.Exec(`INSERT INTO app_metadata (key, value) VALUES ('schema_version', ?)`, schemaVersion); err != nil {
				conn.Close()
				t.Fatalf("failed inserting schema_version: %v", err)
			}
		}
	}

	return conn
}

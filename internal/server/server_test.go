// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/server_test.go

package server

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestFrontendDate(t *testing.T) {
	server := test_server(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/date/2026-04-12", nil)

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if !strings.Contains(response.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("expected HTML content type, got %q", response.Header().Get("Content-Type"))
	}
	if !strings.Contains(response.Body.String(), "Sunday April 12, 2026 / March 30, 2026") {
		t.Fatal("expected dataheader in HTML")
	}
	if !strings.Contains(response.Body.String(), "Saints") {
		t.Fatal("expected Saints in HTML")
	}
	if !strings.Contains(response.Body.String(), "primary") {
		t.Fatal("expected primary marker in HTML")
	}
}

func TestFrontendMissingDate(t *testing.T) {
	server := test_server(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/date/2099-01-01", nil)

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), "Date Not Found") {
		t.Fatal("expected not-found message in HTML")
	}
}

func TestStaticExport(t *testing.T) {
	server := test_server(t)
	outputDir := t.TempDir()

	count, err := server.ExportWeb(outputDir)
	if err != nil {
		t.Fatalf("ExportWeb returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 exported day, got %d", count)
	}

	paths := []string{
		"index.html",
		"assets/style.css",
		"assets/app.js",
		filepath.Join("dates", "2026-04-12", "index.html"),
		filepath.Join("api", "date", "2026-04-12.json"),
		filepath.Join("api", "info.json"),
	}

	for _, path := range paths {
		fullPath := filepath.Join(outputDir, path)
		if _, err := os.Stat(fullPath); err != nil {
			t.Fatalf("expected exported file %s: %v", fullPath, err)
		}
	}
}

func test_server(t *testing.T) *Server {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "orthocal.db")
	conn := test_server_db(t, dbPath)
	t.Cleanup(func() {
		conn.Close()
	})

	server, err := New(conn, Config{
		DatabasePath: dbPath,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	return server
}

func test_server_db(t *testing.T, path string) *sql.DB {
	t.Helper()

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}

	statements := []string{
		`CREATE TABLE app_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL)`,
		`INSERT INTO app_metadata (key, value) VALUES ('schema_version', 'test')`,
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
		`INSERT INTO calendar_days VALUES (
			1,
			'Sunday April 12, 2026 / March 30, 2026',
			'2026-04-12',
			'Sunday',
			'2026-03-30',
			'The Bright Resurrection of Christ, The Pascha of the Lord.',
			'The End of the Great Lent.'
		)`,
		`INSERT INTO saints VALUES (1, 1, 'Venerable John Climacus of Sinai.', '1.gif', 1, 0, '1', 'without a sign')`,
		`INSERT INTO saints VALUES (1, 2, 'St. Osburga of Coventry.', 'o.gif', 0, 1, 'o', 'ordinary/minor')`,
		`INSERT INTO scripture_readings VALUES (1, 1, 'John 20:19-25', 'Vespers, Gospel', '', 'John 20:19-25 - Vespers, Gospel')`,
		`INSERT INTO hymns VALUES (1, 1, 1, 'Troparion', 'V', 'Pascha Troparion', 'Christ is risen from the dead.')`,
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			conn.Close()
			t.Fatalf("failed statement %q: %v", statement, err)
		}
	}

	return conn
}

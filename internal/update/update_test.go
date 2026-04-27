// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/update/update_test.go

package update

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestIsHTTPSource(t *testing.T) {
	tests := []struct {
		source string
		want   bool
	}{
		{source: "http://example.org/orthodox-calendar.db", want: true},
		{source: "https://example.org/orthodox-calendar.db", want: true},
		{source: "/tmp/orthodox-calendar.db", want: false},
		{source: "file:///tmp/orthodox-calendar.db", want: false},
		{source: "ftp://example.org/orthodox-calendar.db", want: false},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			if got := IsHTTPSource(test.source); got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestUpdateDatabaseCreatesBackup(t *testing.T) {
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "source.db")
	target := filepath.Join(tempDir, "target.db")

	create_update_db(t, source, true, "4")
	create_update_db(t, target, true, "4")

	result, err := UpdateDatabase(target, source, false)
	if err != nil {
		t.Fatalf("UpdateDatabase returned error: %v", err)
	}

	if !result.BackupCreated {
		t.Fatal("expected backup to be created")
	}

	if _, err := os.Stat(target + ".bak"); err != nil {
		t.Fatalf("expected backup file: %v", err)
	}

	if _, err := ValidateDatabase(target, false); err != nil {
		t.Fatalf("target failed validation: %v", err)
	}
}

func TestUpdateDatabaseCopiesLocalSource(t *testing.T) {
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "source.db")
	target := filepath.Join(tempDir, "nested", "target.db")

	create_update_db(t, source, true, "4")

	result, err := UpdateDatabase(target, source, false)
	if err != nil {
		t.Fatalf("UpdateDatabase returned error: %v", err)
	}

	if result.Source != source || result.TargetPath != target {
		t.Fatalf("unexpected result: %#v", result)
	}

	if result.BytesWritten <= 0 {
		t.Fatalf("expected bytes written, got %d", result.BytesWritten)
	}

	if result.BackupCreated {
		t.Fatal("did not expect backup")
	}

	if _, err := ValidateDatabase(target, false); err != nil {
		t.Fatalf("target failed validation: %v", err)
	}
}

func TestUpdateDatabaseRejectsSamePath(t *testing.T) {
	tempDir := t.TempDir()
	source := filepath.Join(tempDir, "source.db")
	create_update_db(t, source, true, "4")

	if _, err := UpdateDatabase(source, source, false); err == nil {
		t.Fatal("expected same path error")
	}
}

func TestValidateDatabaseAcceptsValidDB(t *testing.T) {
	path := filepath.Join(t.TempDir(), "valid.db")
	create_update_db(t, path, true, "4")

	if _, err := ValidateDatabase(path, false); err != nil {
		t.Fatalf("ValidateDatabase returned error: %v", err)
	}
}

func TestValidateDatabaseRejectsMissingMetadata(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-metadata.db")
	create_update_db(t, path, false, "")

	if _, err := ValidateDatabase(path, false); err == nil {
		t.Fatal("expected missing metadata error")
	}
}

func TestValidateDatabaseRejectsMissingSchemaVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-schema-version.db")
	create_update_db(t, path, true, "")

	if _, err := ValidateDatabase(path, false); err == nil {
		t.Fatal("expected missing schema_version error")
	}
}

func TestValidateDatabaseRejectsNewerSchemaWithoutForce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "newer-schema.db")
	create_update_db(t, path, true, "5")

	if _, err := ValidateDatabase(path, false); err == nil {
		t.Fatal("expected newer schema error")
	}

	result, err := ValidateDatabase(path, true)
	if err != nil {
		t.Fatalf("ValidateDatabase force returned error: %v", err)
	}

	if !result.Forced || result.SchemaVersion != 5 {
		t.Fatalf("unexpected validation result: %#v", result)
	}
}

func create_update_db(t *testing.T, path string, metadata bool, schemaVersion string) {
	t.Helper()

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("sql.Open returned error: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Exec(`CREATE TABLE calendar_days (id INTEGER PRIMARY KEY)`); err != nil {
		t.Fatalf("failed creating calendar_days: %v", err)
	}

	if metadata {
		if _, err := conn.Exec(`CREATE TABLE app_metadata (key TEXT PRIMARY KEY, value TEXT NOT NULL)`); err != nil {
			t.Fatalf("failed creating app_metadata: %v", err)
		}
	}

	if schemaVersion != "" {
		if _, err := conn.Exec(`INSERT INTO app_metadata (key, value) VALUES ('schema_version', ?)`, schemaVersion); err != nil {
			t.Fatalf("failed inserting schema_version: %v", err)
		}
	}
}

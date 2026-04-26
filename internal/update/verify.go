// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/update/verify.go

package update

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

func ValidateDatabase(path string) error {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		return err
	}

	var integrityResult string
	if err := conn.QueryRow("PRAGMA integrity_check").Scan(&integrityResult); err != nil {
		return err
	}

	if integrityResult != "ok" {
		return fmt.Errorf("database integrity check failed: %s", integrityResult)
	}

	hasMetadata, err := table_exists(conn, "app_metadata")
	if err != nil {
		return err
	}

	if !hasMetadata {
		return errors.New("database is missing app_metadata table")
	}

	var schemaVersion string
	err = conn.QueryRow("SELECT value FROM app_metadata WHERE key = 'schema_version'").Scan(&schemaVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("database is missing app_metadata schema_version")
		}

		return err
	}

	if schemaVersion == "" {
		return errors.New("database has empty app_metadata schema_version")
	}

	return nil
}

func table_exists(conn *sql.DB, table string) (bool, error) {
	var name string
	err := conn.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?", table).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

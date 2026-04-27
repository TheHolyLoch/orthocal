// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/update/verify.go

package update

import (
	"database/sql"
	"errors"
	"fmt"

	"orthocal/internal/db"

	_ "modernc.org/sqlite"
)

type ValidationResult struct {
	SchemaVersion int
	Forced        bool
}

func ValidateDatabase(path string, force bool) (ValidationResult, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return ValidationResult{}, err
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		return ValidationResult{}, err
	}

	var integrityResult string
	if err := conn.QueryRow("PRAGMA integrity_check").Scan(&integrityResult); err != nil {
		return ValidationResult{}, err
	}

	if integrityResult != "ok" {
		return ValidationResult{}, fmt.Errorf("database integrity check failed: %s", integrityResult)
	}

	hasMetadata, err := table_exists(conn, "app_metadata")
	if err != nil {
		return ValidationResult{}, err
	}

	if !hasMetadata {
		return ValidationResult{}, errors.New("database is missing app_metadata table")
	}

	schemaVersion, known, err := db.SchemaVersion(conn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ValidationResult{}, errors.New("database is missing app_metadata schema_version")
		}

		return ValidationResult{}, err
	}

	if !known {
		return ValidationResult{}, errors.New("database is missing app_metadata schema_version")
	}

	result := ValidationResult{
		SchemaVersion: schemaVersion,
	}

	if schemaVersion > db.SupportedSchemaVersion {
		if !force {
			return ValidationResult{}, fmt.Errorf("database schema_version %d is newer than supported schema_version %d; use --force to update anyway", schemaVersion, db.SupportedSchemaVersion)
		}

		result.Forced = true
	}

	return result, nil
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

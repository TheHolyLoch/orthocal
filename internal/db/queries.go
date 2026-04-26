// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/queries.go

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrTableMissing = errors.New("table missing")

func CountRows(conn *sql.DB, table string) (int, error) {
	if !valid_table_name(table) {
		return 0, fmt.Errorf("invalid table name: %s", table)
	}

	exists, err := table_exists(conn, table)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, nil
	}

	var count int
	if err := conn.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func MetadataRows(conn *sql.DB) ([]Metadata, error) {
	exists, err := table_exists(conn, "app_metadata")
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrTableMissing
	}

	rows, err := conn.Query("SELECT key, value FROM app_metadata ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metadata := []Metadata{}
	for rows.Next() {
		item := Metadata{}
		if err := rows.Scan(&item.Key, &item.Value); err != nil {
			return nil, err
		}

		metadata = append(metadata, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metadata, nil
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

func valid_table_name(table string) bool {
	if table == "" {
		return false
	}

	for _, char := range table {
		if char != '_' && (char < '0' || char > '9') && (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
			return false
		}
	}

	return !strings.HasPrefix(table, "sqlite_")
}

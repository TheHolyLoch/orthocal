// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/db.go

package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("database file does not exist: %s", path)
		}

		return nil, err
	}

	conn, err := sql.Open("sqlite", sqlite_read_only_dsn(path))
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func sqlite_read_only_dsn(path string) string {
	if strings.HasPrefix(path, "file:") {
		return path
	}

	absolute, err := filepath.Abs(path)
	if err == nil {
		path = absolute
	}

	return (&url.URL{
		Scheme:   "file",
		Path:     path,
		RawQuery: "mode=ro",
	}).String()
}

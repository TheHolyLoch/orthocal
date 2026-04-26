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

func DayByGregorianDate(conn *sql.DB, value string) (CalendarDay, bool, error) {
	row := conn.QueryRow(`
		SELECT
			id,
			dataheader,
			gregorian_date,
			gregorian_weekday,
			julian_date,
			headerheader,
			fasting_rule
		FROM calendar_days
		WHERE gregorian_date = ?
	`, value)

	day := CalendarDay{}
	if err := row.Scan(
		&day.ID,
		&day.DataHeader,
		&day.GregorianDate,
		&day.GregorianWeekday,
		&day.JulianDate,
		&day.HeaderHeader,
		&day.FastingRule,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CalendarDay{}, false, nil
		}

		return CalendarDay{}, false, err
	}

	return day, true, nil
}

func DayViewByGregorianDate(conn *sql.DB, value string) (DayView, bool, error) {
	day, found, err := DayByGregorianDate(conn, value)
	if err != nil || !found {
		return DayView{}, found, err
	}

	saints, err := SaintsByDayID(conn, day.ID)
	if err != nil {
		return DayView{}, false, err
	}

	scripture, err := ScriptureByDayID(conn, day.ID)
	if err != nil {
		return DayView{}, false, err
	}

	return DayView{
		Day:               day,
		PrimarySaints:     primary_saints(saints),
		WesternSaints:     western_saints(saints),
		ScriptureReadings: scripture,
		Saints:            saints,
	}, true, nil
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

func SaintsByDayID(conn *sql.DB, dayID int) ([]Saint, error) {
	exists, err := table_exists(conn, "saints")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []Saint{}, nil
	}

	rows, err := conn.Query(`
		SELECT
			saint_order,
			name,
			icon_file,
			is_primary,
			is_western,
			service_rank_code,
			service_rank_name
		FROM saints
		WHERE day_id = ?
		ORDER BY saint_order
	`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	saints := []Saint{}
	for rows.Next() {
		item := Saint{}
		isPrimary := 0
		isWestern := 0

		if err := rows.Scan(
			&item.SaintOrder,
			&item.Name,
			&item.IconFile,
			&isPrimary,
			&isWestern,
			&item.ServiceRankCode,
			&item.ServiceRankName,
		); err != nil {
			return nil, err
		}

		item.IsPrimary = isPrimary == 1
		item.IsWestern = isWestern == 1
		saints = append(saints, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return saints, nil
}

func ScriptureByDayID(conn *sql.DB, dayID int) ([]ScriptureReading, error) {
	exists, err := table_exists(conn, "scripture_readings")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []ScriptureReading{}, nil
	}

	rows, err := conn.Query(`
		SELECT
			reading_order,
			verse_reference,
			description,
			reading_url,
			display_text
		FROM scripture_readings
		WHERE day_id = ?
		ORDER BY reading_order
	`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	readings := []ScriptureReading{}
	for rows.Next() {
		item := ScriptureReading{}
		if err := rows.Scan(
			&item.ReadingOrder,
			&item.VerseReference,
			&item.Description,
			&item.ReadingURL,
			&item.DisplayText,
		); err != nil {
			return nil, err
		}

		readings = append(readings, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return readings, nil
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

func primary_saints(saints []Saint) []Saint {
	filtered := []Saint{}
	for _, saint := range saints {
		if saint.IsPrimary {
			filtered = append(filtered, saint)
		}
	}

	return filtered
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

func western_saints(saints []Saint) []Saint {
	filtered := []Saint{}
	for _, saint := range saints {
		if saint.IsWestern {
			filtered = append(filtered, saint)
		}
	}

	return filtered
}

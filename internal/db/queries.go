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

func ClampLimit(limit int) int {
	if limit <= 0 {
		return 25
	}

	if limit > 200 {
		return 200
	}

	return limit
}

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

func EscapeLikeQuery(query string) string {
	query = strings.ReplaceAll(query, `\`, `\\`)
	query = strings.ReplaceAll(query, `%`, `\%`)
	query = strings.ReplaceAll(query, `_`, `\_`)

	return query
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

func HymnsByDayID(conn *sql.DB, dayID int) ([]Hymn, error) {
	exists, err := table_exists(conn, "hymns")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []Hymn{}, nil
	}

	rows, err := conn.Query(`
		SELECT
			hymn_order,
			section_order,
			hymn_type,
			tone,
			title,
			text
		FROM hymns
		WHERE day_id = ?
		ORDER BY section_order, hymn_order
	`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hymns := []Hymn{}
	for rows.Next() {
		item := Hymn{}
		if err := rows.Scan(
			&item.HymnOrder,
			&item.SectionOrder,
			&item.HymnType,
			&item.Tone,
			&item.Title,
			&item.Text,
		); err != nil {
			return nil, err
		}

		hymns = append(hymns, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hymns, nil
}

func HymnsViewByGregorianDate(conn *sql.DB, value string) (HymnsView, bool, error) {
	day, found, err := DayByGregorianDate(conn, value)
	if err != nil || !found {
		return HymnsView{}, found, err
	}

	hymns, err := HymnsByDayID(conn, day.ID)
	if err != nil {
		return HymnsView{}, false, err
	}

	return HymnsView{
		Day:   day,
		Hymns: hymns,
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

func ReadingsViewByGregorianDate(conn *sql.DB, value string) (ReadingsView, bool, error) {
	day, found, err := DayByGregorianDate(conn, value)
	if err != nil || !found {
		return ReadingsView{}, found, err
	}

	readings, err := ScriptureByDayID(conn, day.ID)
	if err != nil {
		return ReadingsView{}, false, err
	}

	return ReadingsView{
		Day:               day,
		ScriptureReadings: readings,
	}, true, nil
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

func SaintsViewByGregorianDate(conn *sql.DB, value string) (SaintsView, bool, error) {
	day, found, err := DayByGregorianDate(conn, value)
	if err != nil || !found {
		return SaintsView{}, found, err
	}

	saints, err := SaintsByDayID(conn, day.ID)
	if err != nil {
		return SaintsView{}, false, err
	}

	return SaintsView{
		Day:    day,
		Saints: saints,
	}, true, nil
}

func SearchHymns(conn *sql.DB, query string, limit int) ([]SearchResultHymn, error) {
	exists, err := table_exists(conn, "hymns")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []SearchResultHymn{}, nil
	}

	likeQuery := "%" + EscapeLikeQuery(query) + "%"
	rows, err := conn.Query(`
		SELECT
			calendar_days.gregorian_date,
			calendar_days.julian_date,
			hymns.section_order,
			hymns.hymn_order,
			hymns.hymn_type,
			hymns.tone,
			hymns.title,
			substr(hymns.text, 1, 160)
		FROM hymns
		JOIN calendar_days ON calendar_days.id = hymns.day_id
		WHERE hymns.title LIKE ? ESCAPE '\'
			OR hymns.hymn_type LIKE ? ESCAPE '\'
			OR hymns.tone LIKE ? ESCAPE '\'
			OR hymns.text LIKE ? ESCAPE '\'
		ORDER BY calendar_days.gregorian_date, hymns.section_order, hymns.hymn_order
		LIMIT ?
	`, likeQuery, likeQuery, likeQuery, likeQuery, ClampLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []SearchResultHymn{}
	for rows.Next() {
		item := SearchResultHymn{}
		if err := rows.Scan(
			&item.GregorianDate,
			&item.JulianDate,
			&item.SectionOrder,
			&item.HymnOrder,
			&item.HymnType,
			&item.Tone,
			&item.Title,
			&item.TextPreview,
		); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func SearchReadings(conn *sql.DB, query string, limit int) ([]SearchResultReading, error) {
	exists, err := table_exists(conn, "scripture_readings")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []SearchResultReading{}, nil
	}

	likeQuery := "%" + EscapeLikeQuery(query) + "%"
	rows, err := conn.Query(`
		SELECT
			calendar_days.gregorian_date,
			calendar_days.julian_date,
			scripture_readings.reading_order,
			scripture_readings.verse_reference,
			scripture_readings.description
		FROM scripture_readings
		JOIN calendar_days ON calendar_days.id = scripture_readings.day_id
		WHERE scripture_readings.verse_reference LIKE ? ESCAPE '\'
			OR scripture_readings.description LIKE ? ESCAPE '\'
			OR scripture_readings.display_text LIKE ? ESCAPE '\'
		ORDER BY calendar_days.gregorian_date, scripture_readings.reading_order
		LIMIT ?
	`, likeQuery, likeQuery, likeQuery, ClampLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []SearchResultReading{}
	for rows.Next() {
		item := SearchResultReading{}
		if err := rows.Scan(
			&item.GregorianDate,
			&item.JulianDate,
			&item.ReadingOrder,
			&item.VerseReference,
			&item.Description,
		); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func SearchSaints(conn *sql.DB, query string, westernOnly bool, primaryOnly bool, limit int) ([]SearchResultSaint, error) {
	exists, err := table_exists(conn, "saints")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []SearchResultSaint{}, nil
	}

	likeQuery := "%" + EscapeLikeQuery(query) + "%"
	rows, err := conn.Query(`
		SELECT
			calendar_days.gregorian_date,
			calendar_days.julian_date,
			saints.saint_order,
			saints.name,
			saints.service_rank_code,
			saints.service_rank_name,
			saints.is_primary,
			saints.is_western
		FROM saints
		JOIN calendar_days ON calendar_days.id = saints.day_id
		WHERE saints.name LIKE ? ESCAPE '\'
			AND (? = 0 OR saints.is_western = 1)
			AND (? = 0 OR saints.is_primary = 1)
		ORDER BY calendar_days.gregorian_date, saints.saint_order
		LIMIT ?
	`, likeQuery, bool_int(westernOnly), bool_int(primaryOnly), ClampLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []SearchResultSaint{}
	for rows.Next() {
		item := SearchResultSaint{}
		isPrimary := 0
		isWestern := 0
		if err := rows.Scan(
			&item.GregorianDate,
			&item.JulianDate,
			&item.SaintOrder,
			&item.Name,
			&item.ServiceRankCode,
			&item.ServiceRankName,
			&isPrimary,
			&isWestern,
		); err != nil {
			return nil, err
		}

		item.IsPrimary = isPrimary == 1
		item.IsWestern = isWestern == 1
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
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

func bool_int(value bool) int {
	if value {
		return 1
	}

	return 0
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

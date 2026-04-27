// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/queries.go

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var ErrTableMissing = errors.New("table missing")

const SupportedSchemaVersion = 4

func AllCalendarDates(conn *sql.DB) ([]string, error) {
	rows, err := conn.Query("SELECT gregorian_date FROM calendar_days ORDER BY gregorian_date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dates := []string{}
	for rows.Next() {
		value := ""
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}

		dates = append(dates, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dates, nil
}

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

func SchemaCompatibilityStatus(conn *sql.DB) (CompatibilityStatus, error) {
	schemaVersion, known, err := SchemaVersion(conn)
	if err != nil {
		return CompatibilityStatus{}, err
	}

	status := CompatibilityStatus{
		SchemaVersion:    schemaVersion,
		SchemaKnown:      known,
		SupportedVersion: SupportedSchemaVersion,
	}

	switch {
	case !known:
		status.Message = "schema_version is unknown"

	case schemaVersion < SupportedSchemaVersion:
		status.IsOlder = true
		status.Message = "older database schema; calculated calendar events may be unavailable"

	case schemaVersion > SupportedSchemaVersion:
		status.IsNewer = true
		status.Message = "newer database schema; read-only commands will try compatible queries"

	default:
		status.Message = "schema is compatible"
	}

	return status, nil
}

func EscapeLikeQuery(query string) string {
	query = strings.ReplaceAll(query, `\`, `\\`)
	query = strings.ReplaceAll(query, `%`, `\%`)
	query = strings.ReplaceAll(query, `_`, `\_`)

	return query
}

func FirstCalendarDate(conn *sql.DB) (string, bool, error) {
	value := ""
	err := conn.QueryRow("SELECT gregorian_date FROM calendar_days ORDER BY gregorian_date LIMIT 1").Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil
		}

		return "", false, err
	}

	return value, true, nil
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

	hasV4Columns, err := column_exists(conn, "calendar_days", "feasts")
	if err != nil {
		return CalendarDay{}, false, err
	}

	if hasV4Columns {
		if err := scan_day_v4(conn, &day); err != nil {
			return CalendarDay{}, false, err
		}
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

	events, err := CalendarDayEventsByDayID(conn, day.ID)
	if err != nil {
		return DayView{}, false, err
	}

	return DayView{
		Day:               day,
		Events:            events,
		FeastEvents:       feast_events(events),
		FastEvents:        fast_events(events),
		RemembranceEvents: remembrance_events(events),
		FastFreeEvents:    fast_free_events(events),
		PrimarySaints:     primary_saints(saints),
		WesternSaints:     western_saints(saints),
		ScriptureReadings: scripture,
		Saints:            saints,
	}, true, nil
}

func CalendarDayEventsByDayID(conn *sql.DB, dayID int) ([]CalendarDayEvent, error) {
	exists, err := table_exists(conn, "calendar_day_events")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []CalendarDayEvent{}, nil
	}

	rows, err := conn.Query(`
		SELECT
			day_id,
			event_id,
			event_date,
			category,
			title,
			sort_order
		FROM calendar_day_events
		WHERE day_id = ?
		ORDER BY sort_order, title
	`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []CalendarDayEvent{}
	for rows.Next() {
		item := CalendarDayEvent{}
		if err := rows.Scan(
			&item.DayID,
			&item.EventID,
			&item.EventDate,
			&item.Category,
			&item.Title,
			&item.SortOrder,
		); err != nil {
			return nil, err
		}

		events = append(events, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
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

func InfoViewByPath(conn *sql.DB, path string) (InfoView, error) {
	view := InfoView{
		DatabasePath: path,
		Metadata:     []Metadata{},
	}

	metadata, err := MetadataRows(conn)
	if err != nil {
		if errors.Is(err, ErrTableMissing) {
			view.MetadataUnavailable = true
		} else {
			return InfoView{}, err
		}
	} else {
		view.Metadata = metadata
	}

	compatibility, err := SchemaCompatibilityStatus(conn)
	if err != nil {
		return InfoView{}, err
	}
	view.Compatibility = compatibility

	if compatibility.IsOlder {
		view.SchemaNote = compatibility.Message
	}

	if count, err := CountRows(conn, "calendar_days"); err == nil {
		view.Counts.CalendarDays = count
	} else {
		return InfoView{}, err
	}

	if count, err := CountRows(conn, "calendar_events"); err == nil {
		view.Counts.CalendarEvents = count
	} else {
		return InfoView{}, err
	}

	if count, err := CountRows(conn, "calendar_day_events"); err == nil {
		view.Counts.CalendarDayEvents = count
	} else {
		return InfoView{}, err
	}

	if count, err := CountRows(conn, "saints"); err == nil {
		view.Counts.Saints = count
	} else {
		return InfoView{}, err
	}

	if count, err := CountRows(conn, "scripture_readings"); err == nil {
		view.Counts.ScriptureReadings = count
	} else {
		return InfoView{}, err
	}

	if count, err := CountRows(conn, "hymns"); err == nil {
		view.Counts.Hymns = count
	} else {
		return InfoView{}, err
	}

	return view, nil
}

func SchemaVersion(conn *sql.DB) (int, bool, error) {
	exists, err := table_exists(conn, "app_metadata")
	if err != nil {
		return 0, false, err
	}

	if !exists {
		return 0, false, nil
	}

	value := ""
	err = conn.QueryRow("SELECT value FROM app_metadata WHERE key = 'schema_version'").Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}

		return 0, false, err
	}

	version, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, false, fmt.Errorf("invalid schema_version: %s", value)
	}

	return version, true, nil
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

	sort_saints_for_day(saints)
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

func SearchEvents(conn *sql.DB, query string, category string, limit int) ([]SearchResultEvent, error) {
	exists, err := table_exists(conn, "calendar_events")
	if err != nil {
		return nil, err
	}

	if !exists {
		return []SearchResultEvent{}, nil
	}

	likeQuery := "%" + EscapeLikeQuery(query) + "%"
	categoryFilter := event_search_category(category)
	rows, err := conn.Query(`
		SELECT
			category,
			title,
			start_date,
			end_date,
			is_range
		FROM calendar_events
		WHERE (title LIKE ? ESCAPE '\'
			OR category LIKE ? ESCAPE '\'
			OR notes LIKE ? ESCAPE '\')
			AND (? = ''
				OR (? = 'feast' AND category IN ('fixed_feast', 'movable_feast'))
				OR category = ?)
		ORDER BY start_date, sort_order, title
		LIMIT ?
	`, likeQuery, likeQuery, likeQuery, categoryFilter, categoryFilter, categoryFilter, ClampLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []SearchResultEvent{}
	for rows.Next() {
		item := SearchResultEvent{}
		isRange := 0
		if err := rows.Scan(
			&item.Category,
			&item.Title,
			&item.StartDate,
			&item.EndDate,
			&isRange,
		); err != nil {
			return nil, err
		}

		item.IsRange = isRange == 1
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
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

func column_exists(conn *sql.DB, table string, column string) (bool, error) {
	rows, err := conn.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		cid := 0
		name := ""
		columnType := ""
		notNull := 0
		defaultValue := sql.NullString{}
		primaryKey := 0

		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}

		if name == column {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func event_search_category(category string) string {
	switch category {
	case "feasts":
		return "feast"
	case "fasts":
		return "fasting_season"
	case "remembrances":
		return "remembrance"
	}

	return ""
}

func fast_events(events []CalendarDayEvent) []CalendarDayEvent {
	filtered := []CalendarDayEvent{}
	for _, event := range events {
		if event.Category == "fasting_season" {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

func fast_free_events(events []CalendarDayEvent) []CalendarDayEvent {
	filtered := []CalendarDayEvent{}
	for _, event := range events {
		if event.Category == "fast_free_week" {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

func feast_events(events []CalendarDayEvent) []CalendarDayEvent {
	filtered := []CalendarDayEvent{}
	for _, event := range events {
		if event.Category == "fixed_feast" || event.Category == "movable_feast" {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

func metadata_value(metadata []Metadata, key string) string {
	for _, item := range metadata {
		if item.Key == key {
			return item.Value
		}
	}

	return ""
}

func remembrance_events(events []CalendarDayEvent) []CalendarDayEvent {
	filtered := []CalendarDayEvent{}
	for _, event := range events {
		if event.Category == "remembrance" {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

func scan_day_v4(conn *sql.DB, day *CalendarDay) error {
	isHoliday := 0
	isLentDay := 0
	if err := conn.QueryRow(`
		SELECT
			feasts,
			fasts,
			remembrances,
			fast_free_periods,
			fasting_level_code,
			fasting_level_name,
			fasting_level_description,
			is_holiday,
			is_lent_day
		FROM calendar_days
		WHERE id = ?
	`, day.ID).Scan(
		&day.Feasts,
		&day.Fasts,
		&day.Remembrances,
		&day.FastFreePeriods,
		&day.FastingLevelCode,
		&day.FastingLevelName,
		&day.FastingLevelDescription,
		&isHoliday,
		&isLentDay,
	); err != nil {
		return err
	}

	day.IsHoliday = isHoliday == 1
	day.IsLentDay = isLentDay == 1
	return nil
}

func saint_day_sort_key(saint Saint) (int, int) {
	if saint.IsPrimary {
		return 0, 0
	}

	if saint.IsWestern {
		return 1, 0
	}

	rank, err := strconv.Atoi(saint.ServiceRankCode)
	if err == nil && rank >= 0 && rank <= 6 {
		return 2, rank
	}

	return 3, -1
}

func sort_saints_for_day(saints []Saint) {
	sort.SliceStable(saints, func(left int, right int) bool {
		leftGroup, leftRank := saint_day_sort_key(saints[left])
		rightGroup, rightRank := saint_day_sort_key(saints[right])

		if leftGroup != rightGroup {
			return leftGroup < rightGroup
		}

		if leftRank != rightRank {
			return leftRank > rightRank
		}

		return saints[left].SaintOrder < saints[right].SaintOrder
	})
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

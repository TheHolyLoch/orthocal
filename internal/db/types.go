// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/types.go

package db

type CalendarDay struct {
	ID               int    `json:"id"`
	DataHeader       string `json:"dataheader"`
	GregorianDate    string `json:"gregorian_date"`
	GregorianWeekday string `json:"gregorian_weekday"`
	JulianDate       string `json:"julian_date"`
	HeaderHeader     string `json:"headerheader"`
	FastingRule      string `json:"fasting_rule"`
}

type DayView struct {
	Day               CalendarDay        `json:"day"`
	PrimarySaints     []Saint            `json:"primary_saints"`
	WesternSaints     []Saint            `json:"western_saints"`
	ScriptureReadings []ScriptureReading `json:"scripture_readings"`
	Saints            []Saint            `json:"saints"`
}

type Hymn struct{}

type InfoCounts struct {
	CalendarDays      int `json:"calendar_days"`
	Saints            int `json:"saints"`
	ScriptureReadings int `json:"scripture_readings"`
	Hymns             int `json:"hymns"`
}

type InfoView struct {
	DatabasePath        string     `json:"database_path"`
	Metadata            []Metadata `json:"metadata"`
	Counts              InfoCounts `json:"counts"`
	MetadataUnavailable bool       `json:"-"`
}

type Metadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Saint struct {
	SaintOrder      int    `json:"saint_order"`
	Name            string `json:"name"`
	IconFile        string `json:"icon_file"`
	IsPrimary       bool   `json:"is_primary"`
	IsWestern       bool   `json:"is_western"`
	ServiceRankCode string `json:"service_rank_code"`
	ServiceRankName string `json:"service_rank_name"`
}

type ScriptureReading struct {
	ReadingOrder   int    `json:"reading_order"`
	VerseReference string `json:"verse_reference"`
	Description    string `json:"description"`
	ReadingURL     string `json:"reading_url"`
	DisplayText    string `json:"display_text"`
}

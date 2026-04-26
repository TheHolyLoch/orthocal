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

type Hymn struct {
	HymnOrder    int    `json:"hymn_order"`
	SectionOrder int    `json:"section_order"`
	HymnType     string `json:"hymn_type"`
	Tone         string `json:"tone"`
	Title        string `json:"title"`
	Text         string `json:"text"`
}

type HymnsView struct {
	Day   CalendarDay `json:"day"`
	Hymns []Hymn      `json:"hymns"`
}

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

type SaintsView struct {
	Day    CalendarDay `json:"day"`
	Saints []Saint     `json:"saints"`
}

type ScriptureReading struct {
	ReadingOrder   int    `json:"reading_order"`
	VerseReference string `json:"verse_reference"`
	Description    string `json:"description"`
	ReadingURL     string `json:"reading_url"`
	DisplayText    string `json:"display_text"`
}

type ReadingsView struct {
	Day               CalendarDay        `json:"day"`
	ScriptureReadings []ScriptureReading `json:"scripture_readings"`
}

type SearchHymnsView struct {
	Query    string             `json:"query"`
	Category string             `json:"category"`
	Limit    int                `json:"limit"`
	Results  []SearchResultHymn `json:"results"`
}

type SearchReadingsView struct {
	Query    string                `json:"query"`
	Category string                `json:"category"`
	Limit    int                   `json:"limit"`
	Results  []SearchResultReading `json:"results"`
}

type SearchResultHymn struct {
	GregorianDate string `json:"gregorian_date"`
	JulianDate    string `json:"julian_date"`
	SectionOrder  int    `json:"section_order"`
	HymnOrder     int    `json:"hymn_order"`
	HymnType      string `json:"hymn_type"`
	Tone          string `json:"tone"`
	Title         string `json:"title"`
	TextPreview   string `json:"text_preview"`
}

type SearchResultReading struct {
	GregorianDate  string `json:"gregorian_date"`
	JulianDate     string `json:"julian_date"`
	ReadingOrder   int    `json:"reading_order"`
	VerseReference string `json:"verse_reference"`
	Description    string `json:"description"`
}

type SearchResultSaint struct {
	GregorianDate   string `json:"gregorian_date"`
	JulianDate      string `json:"julian_date"`
	SaintOrder      int    `json:"saint_order"`
	Name            string `json:"name"`
	ServiceRankCode string `json:"service_rank_code"`
	ServiceRankName string `json:"service_rank_name"`
	IsPrimary       bool   `json:"is_primary"`
	IsWestern       bool   `json:"is_western"`
}

type SearchSaintsView struct {
	Query    string              `json:"query"`
	Category string              `json:"category"`
	Limit    int                 `json:"limit"`
	Results  []SearchResultSaint `json:"results"`
}

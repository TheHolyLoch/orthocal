// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/types.go

package db

type CalendarDay struct {
	ID               int
	DataHeader       string
	GregorianDate    string
	GregorianWeekday string
	JulianDate       string
	HeaderHeader     string
	FastingRule      string
}

type DayView struct {
	Day               CalendarDay
	Saints            []Saint
	ScriptureReadings []ScriptureReading
}

type Hymn struct{}

type InfoView struct {
	DatabasePath           string
	Metadata               []Metadata
	MetadataUnavailable    bool
	CalendarDaysCount      int
	SaintsCount            int
	ScriptureReadingsCount int
	HymnsCount             int
}

type Metadata struct {
	Key   string
	Value string
}

type Saint struct {
	SaintOrder      int
	Name            string
	IconFile        string
	IsPrimary       bool
	IsWestern       bool
	ServiceRankCode string
	ServiceRankName string
}

type ScriptureReading struct {
	ReadingOrder   int
	VerseReference string
	Description    string
	ReadingURL     string
	DisplayText    string
}

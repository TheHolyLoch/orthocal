// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/db/types.go

package db

type CalendarDay struct{}

type DayView struct{}

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

type Saint struct{}

type ScriptureReading struct{}

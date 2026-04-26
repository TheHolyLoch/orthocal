// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/render/render.go

package render

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"orthocal/internal/db"
)

func Day(output io.Writer, view db.DayView) {
	fmt.Fprintf(output, "%s / %s\n", format_gregorian_date(view.Day), format_julian_date(view.Day.JulianDate))

	if strings.TrimSpace(view.Day.HeaderHeader) != "" {
		fmt.Fprintf(output, "%s\n", view.Day.HeaderHeader)
	}

	if strings.TrimSpace(view.Day.FastingRule) != "" {
		fmt.Fprintf(output, "%s\n", view.Day.FastingRule)
	}

	primarySaints := primary_saints(view.Saints)
	if len(primarySaints) > 0 {
		fmt.Fprintln(output)
		fmt.Fprintln(output, "Primary saints:")
		for index, saint := range primarySaints {
			fmt.Fprintf(output, "\t%d. %s\n", index+1, saint.Name)
		}
	}

	westernSaints := western_saints(view.Saints)
	if len(westernSaints) > 0 {
		fmt.Fprintln(output)
		fmt.Fprintln(output, "Western saints:")
		for _, saint := range westernSaints {
			fmt.Fprintf(output, "\t- %s\n", saint.Name)
		}
	}

	if len(view.ScriptureReadings) > 0 {
		fmt.Fprintln(output)
		fmt.Fprintln(output, "Scripture:")
		for _, reading := range view.ScriptureReadings {
			if strings.TrimSpace(reading.Description) != "" {
				fmt.Fprintf(output, "\t- %s - %s\n", reading.VerseReference, reading.Description)
			} else {
				fmt.Fprintf(output, "\t- %s\n", reading.VerseReference)
			}
		}
	}
}

func Hymns(output io.Writer, view db.HymnsView) {
	fmt.Fprintf(output, "%s / %s\n", format_gregorian_date(view.Day), format_julian_date(view.Day.JulianDate))
	fmt.Fprintln(output)

	if len(view.Hymns) == 0 {
		fmt.Fprintln(output, "No hymns found.")
		return
	}

	fmt.Fprintln(output, "Hymns:")
	for _, hymn := range view.Hymns {
		label := hymn_label(hymn)
		if label != "" {
			fmt.Fprintf(output, "\t%d. %s\n", hymn.HymnOrder, label)
		} else {
			fmt.Fprintf(output, "\t%d.\n", hymn.HymnOrder)
		}

		if strings.TrimSpace(hymn.Text) != "" {
			fmt.Fprintf(output, "\t%s\n", hymn.Text)
		}
	}
}

func Readings(output io.Writer, view db.ReadingsView) {
	fmt.Fprintf(output, "%s / %s\n", format_gregorian_date(view.Day), format_julian_date(view.Day.JulianDate))
	fmt.Fprintln(output)

	if len(view.ScriptureReadings) == 0 {
		fmt.Fprintln(output, "No scripture readings found.")
		return
	}

	fmt.Fprintln(output, "Scripture:")
	for _, reading := range view.ScriptureReadings {
		if strings.TrimSpace(reading.Description) != "" {
			fmt.Fprintf(output, "\t%d. %s - %s\n", reading.ReadingOrder, reading.VerseReference, reading.Description)
		} else {
			fmt.Fprintf(output, "\t%d. %s\n", reading.ReadingOrder, reading.VerseReference)
		}
	}
}

func RenderDayJSON(view db.DayView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderHymnsJSON(view db.HymnsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderInfoJSON(view db.InfoView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderReadingsJSON(view db.ReadingsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderSaintsJSON(view db.SaintsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func Saints(output io.Writer, view db.SaintsView) {
	fmt.Fprintf(output, "%s / %s\n", format_gregorian_date(view.Day), format_julian_date(view.Day.JulianDate))
	fmt.Fprintln(output)

	if len(view.Saints) == 0 {
		fmt.Fprintln(output, "No saints found.")
		return
	}

	fmt.Fprintln(output, "Saints:")
	for _, saint := range view.Saints {
		fmt.Fprintf(output, "\t%d. %s%s\n", saint.SaintOrder, saint_prefix(saint), saint.Name)
	}
}

func format_gregorian_date(day db.CalendarDay) string {
	parsed, err := time.Parse("2006-01-02", day.GregorianDate)
	if err != nil {
		return day.GregorianDate
	}

	return parsed.Format("Monday January 2, 2006")
}

func hymn_label(hymn db.Hymn) string {
	parts := []string{}

	if strings.TrimSpace(hymn.Title) != "" {
		parts = append(parts, hymn.Title)
	}

	if strings.TrimSpace(hymn.HymnType) != "" {
		parts = append(parts, hymn.HymnType)
	}

	if strings.TrimSpace(hymn.Tone) != "" {
		parts = append(parts, hymn.Tone)
	}

	return strings.Join(parts, " - ")
}

func format_julian_date(value string) string {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return value
	}

	return parsed.Format("January 2, 2006")
}

func primary_saints(saints []db.Saint) []db.Saint {
	filtered := []db.Saint{}
	for _, saint := range saints {
		if saint.IsPrimary {
			filtered = append(filtered, saint)
		}
	}

	return filtered
}

func saint_prefix(saint db.Saint) string {
	parts := []string{}

	if strings.TrimSpace(saint.ServiceRankCode) != "" || strings.TrimSpace(saint.ServiceRankName) != "" {
		parts = append(parts, fmt.Sprintf("[%s: %s]", saint.ServiceRankCode, saint.ServiceRankName))
	}

	if saint.IsPrimary {
		parts = append(parts, "[primary]")
	}

	if saint.IsWestern {
		parts = append(parts, "[western]")
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " ") + " "
}

func western_saints(saints []db.Saint) []db.Saint {
	filtered := []db.Saint{}
	for _, saint := range saints {
		if saint.IsWestern {
			filtered = append(filtered, saint)
		}
	}

	return filtered
}

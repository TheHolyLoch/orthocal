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

	if strings.TrimSpace(view.Day.FastingLevelName) != "" {
		fmt.Fprintf(output, "Calculated: %s", view.Day.FastingLevelName)
		if strings.TrimSpace(view.Day.FastingLevelCode) != "" {
			fmt.Fprintf(output, " (%s)", view.Day.FastingLevelCode)
		}
		fmt.Fprintln(output)
	}

	render_event_section(output, "Feasts:", event_titles(view.FeastEvents, view.Day.Feasts))
	render_event_section(output, "Fasting season:", event_titles(view.FastEvents, view.Day.Fasts))
	render_event_section(output, "Fast-free period:", event_titles(view.FastFreeEvents, view.Day.FastFreePeriods))
	render_event_section(output, "Remembrance:", event_titles(view.RemembranceEvents, view.Day.Remembrances))

	saints := view.Saints
	if len(saints) > 0 {
		fmt.Fprintln(output)
		fmt.Fprintln(output, "Saints:")
		for index, saint := range saints {
			fmt.Fprintf(output, "\t%d. %s%s\n", index+1, saint_prefix(saint), saint.Name)
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

func Events(output io.Writer, view db.EventsView) {
	fmt.Fprintf(output, "%s / %s\n", format_gregorian_date(view.Day), format_julian_date(view.Day.JulianDate))
	fmt.Fprintln(output)

	if len(view.Events) == 0 {
		fmt.Fprintln(output, "No calendar events found.")
		return
	}

	fmt.Fprintln(output, "Events:")
	for _, event := range view.Events {
		fmt.Fprintf(output, "\t- %s: %s\n", event.Category, event.Title)
	}
}

func RenderEventsJSON(view db.EventsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
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

func RenderSearchHymnsJSON(view db.SearchHymnsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderSearchEventsJSON(view db.SearchEventsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderSearchReadingsJSON(view db.SearchReadingsView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderSearchSaintsJSON(view db.SearchSaintsView) error {
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

func SearchHymns(output io.Writer, view db.SearchHymnsView) {
	if len(view.Results) == 0 {
		fmt.Fprintln(output, "No hymn results found.")
		return
	}

	for _, result := range view.Results {
		label := search_hymn_label(result)
		if label != "" {
			fmt.Fprintf(output, "%s / %s  %s - %s\n", result.GregorianDate, result.JulianDate, label, result.TextPreview)
		} else {
			fmt.Fprintf(output, "%s / %s  %s\n", result.GregorianDate, result.JulianDate, result.TextPreview)
		}
	}
}

func SearchEvents(output io.Writer, view db.SearchEventsView) {
	if len(view.Results) == 0 {
		fmt.Fprintln(output, "No event results found.")
		return
	}

	for _, result := range view.Results {
		rangeMarker := ""
		if result.IsRange {
			rangeMarker = " - " + result.EndDate
		}
		fmt.Fprintf(output, "%s%s  %s  %s\n", result.StartDate, rangeMarker, result.Category, result.Title)
	}
}

func SearchReadings(output io.Writer, view db.SearchReadingsView) {
	if len(view.Results) == 0 {
		fmt.Fprintln(output, "No scripture reading results found.")
		return
	}

	for _, result := range view.Results {
		if strings.TrimSpace(result.Description) != "" {
			fmt.Fprintf(output, "%s / %s  %s - %s\n", result.GregorianDate, result.JulianDate, result.VerseReference, result.Description)
		} else {
			fmt.Fprintf(output, "%s / %s  %s\n", result.GregorianDate, result.JulianDate, result.VerseReference)
		}
	}
}

func SearchSaints(output io.Writer, view db.SearchSaintsView) {
	if len(view.Results) == 0 {
		fmt.Fprintln(output, "No saint results found.")
		return
	}

	for _, result := range view.Results {
		fmt.Fprintf(output, "%s / %s  %s%s\n", result.GregorianDate, result.JulianDate, result.Name, search_saint_markers(result))
	}
}

func event_titles(events []db.CalendarDayEvent, fallback string) []string {
	titles := []string{}
	seen := map[string]bool{}

	for _, event := range events {
		title := strings.TrimSpace(event.Title)
		if title != "" && !seen[title] {
			titles = append(titles, title)
			seen[title] = true
		}
	}

	for _, title := range split_pipe(fallback) {
		if !seen[title] {
			titles = append(titles, title)
			seen[title] = true
		}
	}

	return titles
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

func search_hymn_label(result db.SearchResultHymn) string {
	parts := []string{}

	if strings.TrimSpace(result.Title) != "" {
		parts = append(parts, result.Title)
	}

	if strings.TrimSpace(result.HymnType) != "" {
		parts = append(parts, result.HymnType)
	}

	if strings.TrimSpace(result.Tone) != "" {
		parts = append(parts, result.Tone)
	}

	return strings.Join(parts, " ")
}

func search_saint_markers(result db.SearchResultSaint) string {
	parts := []string{}

	if result.IsPrimary {
		parts = append(parts, "[primary]")
	}

	if result.IsWestern {
		parts = append(parts, "[UK+IE]")
	}

	if len(parts) == 0 {
		return ""
	}

	return " " + strings.Join(parts, " ")
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

	if saint.IsPrimary {
		parts = append(parts, "[primary]")
	}

	if saint.IsWestern {
		parts = append(parts, "[UK+IE]")
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " ") + " "
}

func render_event_section(output io.Writer, title string, items []string) {
	if len(items) == 0 {
		return
	}

	fmt.Fprintln(output)
	fmt.Fprintln(output, title)
	for _, item := range items {
		fmt.Fprintf(output, "\t- %s\n", item)
	}
}

func split_pipe(value string) []string {
	parts := []string{}
	for _, item := range strings.Split(value, "|") {
		item = strings.TrimSpace(item)
		if item != "" {
			parts = append(parts, item)
		}
	}

	return parts
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

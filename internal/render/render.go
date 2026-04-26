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

func RenderDayJSON(view db.DayView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func RenderInfoJSON(view db.InfoView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	return encoder.Encode(view)
}

func format_gregorian_date(day db.CalendarDay) string {
	parsed, err := time.Parse("2006-01-02", day.GregorianDate)
	if err != nil {
		return day.GregorianDate
	}

	return parsed.Format("Monday January 2, 2006")
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

func western_saints(saints []db.Saint) []db.Saint {
	filtered := []db.Saint{}
	for _, saint := range saints {
		if saint.IsWestern {
			filtered = append(filtered, saint)
		}
	}

	return filtered
}

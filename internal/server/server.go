// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/server.go

package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orthocal/internal/db"
	"orthocal/internal/version"
)

func New(conn *sql.DB, config Config) (*Server, error) {
	parsed, err := template.New("orthocal").Parse(templates)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:      conn,
		config:    config,
		templates: parsed,
	}, nil
}

func (server *Server) Serve(addr string) error {
	return http.ListenAndServe(addr, server.Handler())
}

func (server *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	server.routes(mux)

	return mux
}

func (server *Server) ExportWeb(outputDir string) (int, error) {
	if info, err := os.Stat(outputDir); err == nil {
		if !info.IsDir() {
			return 0, errors.New("output path is an existing regular file")
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return 0, err
		}
	} else {
		return 0, err
	}

	dates, err := db.AllCalendarDates(server.conn)
	if err != nil {
		return 0, err
	}

	if len(dates) == 0 {
		return 0, errors.New("database has no calendar days")
	}

	indexDate := today()
	if _, found, err := db.DayByGregorianDate(server.conn, indexDate); err != nil {
		return 0, err
	} else if !found {
		indexDate = dates[0]
	}

	if err := write_file(filepath.Join(outputDir, "assets", "style.css"), []byte(stylesheet)); err != nil {
		return 0, err
	}
	if err := write_file(filepath.Join(outputDir, "assets", "app.js"), []byte(javascript)); err != nil {
		return 0, err
	}

	infoView, err := db.InfoViewByPath(server.conn, server.config.DatabasePath)
	if err != nil {
		return 0, err
	}
	infoView.Version = version.Current().Version
	if err := write_json_file(filepath.Join(outputDir, "api", "info.json"), infoView); err != nil {
		return 0, err
	}

	for _, value := range dates {
		if err := server.export_day(outputDir, value); err != nil {
			return 0, err
		}
		if err := server.export_saints(outputDir, value); err != nil {
			return 0, err
		}
		if err := server.export_readings(outputDir, value); err != nil {
			return 0, err
		}
		if err := server.export_hymns(outputDir, value); err != nil {
			return 0, err
		}

		view, found, err := db.DayViewByGregorianDate(server.conn, value)
		if err != nil {
			return 0, err
		}
		if !found {
			continue
		}
		if err := write_json_file(filepath.Join(outputDir, "api", "date", value+".json"), view); err != nil {
			return 0, err
		}
	}

	indexBytes, err := server.render_day_bytes(indexDate, "./")
	if err != nil {
		return 0, err
	}
	if err := write_file(filepath.Join(outputDir, "index.html"), indexBytes); err != nil {
		return 0, err
	}

	return len(dates), nil
}

func (server *Server) export_day(outputDir string, value string) error {
	bytes, err := server.render_day_bytes(value, "../../")
	if err != nil {
		return err
	}

	return write_file(filepath.Join(outputDir, "dates", value, "index.html"), bytes)
}

func (server *Server) export_hymns(outputDir string, value string) error {
	view, found, err := db.HymnsViewByGregorianDate(server.conn, value)
	if err != nil {
		return err
	}
	if !found {
		return errDateNotFound
	}

	bytes, err := server.render_template_bytes(PageData{
		Active:    "hymns",
		DateValue: value,
		HymnsView: view,
		Title:     "Hymns",
	}, "../../")
	if err != nil {
		return err
	}

	return write_file(filepath.Join(outputDir, "hymns", value, "index.html"), bytes)
}

func (server *Server) export_readings(outputDir string, value string) error {
	view, found, err := db.ReadingsViewByGregorianDate(server.conn, value)
	if err != nil {
		return err
	}
	if !found {
		return errDateNotFound
	}

	bytes, err := server.render_template_bytes(PageData{
		Active:       "readings",
		DateValue:    value,
		ReadingsView: view,
		Title:        "Readings",
	}, "../../")
	if err != nil {
		return err
	}

	return write_file(filepath.Join(outputDir, "readings", value, "index.html"), bytes)
}

func (server *Server) export_saints(outputDir string, value string) error {
	view, found, err := db.SaintsViewByGregorianDate(server.conn, value)
	if err != nil {
		return err
	}
	if !found {
		return errDateNotFound
	}

	bytes, err := server.render_template_bytes(PageData{
		Active:     "saints",
		DateValue:  value,
		SaintsView: view,
		Title:      "Saints",
	}, "../../")
	if err != nil {
		return err
	}

	return write_file(filepath.Join(outputDir, "saints", value, "index.html"), bytes)
}

func (server *Server) routes(mux *http.ServeMux) {
	mux.HandleFunc("/", server.handle_index)
	mux.HandleFunc("/date/", server.handle_day)
	mux.HandleFunc("/saints/", server.handle_saints)
	mux.HandleFunc("/readings/", server.handle_readings)
	mux.HandleFunc("/hymns/", server.handle_hymns)
	mux.HandleFunc("/api/date/", server.handle_api_date)
	mux.HandleFunc("/api/today", server.handle_api_today)
	mux.HandleFunc("/api/tomorrow", server.handle_api_tomorrow)
	mux.HandleFunc("/api/saints/", server.handle_api_saints)
	mux.HandleFunc("/api/readings/", server.handle_api_readings)
	mux.HandleFunc("/api/hymns/", server.handle_api_hymns)
	mux.HandleFunc("/api/info", server.handle_api_info)
	mux.HandleFunc("/assets/style.css", server.handle_stylesheet)
	mux.HandleFunc("/assets/app.js", server.handle_javascript)
}

func (server *Server) handle_day(response http.ResponseWriter, request *http.Request) {
	value, status, message, ok := route_date(request, "/date/")
	if !ok {
		if status == http.StatusNotFound {
			http.NotFound(response, request)
		} else {
			server.render_error(response, status, message)
		}
		return
	}

	server.render_day(response, value)
}

func (server *Server) handle_hymns(response http.ResponseWriter, request *http.Request) {
	value, status, message, ok := route_date(request, "/hymns/")
	if !ok {
		if status == http.StatusNotFound {
			http.NotFound(response, request)
		} else {
			server.render_error(response, status, message)
		}
		return
	}

	view, found, err := db.HymnsViewByGregorianDate(server.conn, value)
	if err != nil {
		server.render_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		server.render_not_found(response, value)
		return
	}

	server.render(response, PageData{
		Active:    "hymns",
		DateValue: value,
		HymnsView: view,
		Title:     "Hymns",
		Today:     today(),
	})
}

func (server *Server) handle_index(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(response, request)
		return
	}

	server.render_day(response, today())
}

func (server *Server) handle_javascript(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	fmt.Fprint(response, javascript)
}

func (server *Server) handle_readings(response http.ResponseWriter, request *http.Request) {
	value, status, message, ok := route_date(request, "/readings/")
	if !ok {
		if status == http.StatusNotFound {
			http.NotFound(response, request)
		} else {
			server.render_error(response, status, message)
		}
		return
	}

	view, found, err := db.ReadingsViewByGregorianDate(server.conn, value)
	if err != nil {
		server.render_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		server.render_not_found(response, value)
		return
	}

	server.render(response, PageData{
		Active:       "readings",
		DateValue:    value,
		ReadingsView: view,
		Title:        "Readings",
		Today:        today(),
	})
}

func (server *Server) handle_saints(response http.ResponseWriter, request *http.Request) {
	value, status, message, ok := route_date(request, "/saints/")
	if !ok {
		if status == http.StatusNotFound {
			http.NotFound(response, request)
		} else {
			server.render_error(response, status, message)
		}
		return
	}

	view, found, err := db.SaintsViewByGregorianDate(server.conn, value)
	if err != nil {
		server.render_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		server.render_not_found(response, value)
		return
	}

	server.render(response, PageData{
		Active:     "saints",
		DateValue:  value,
		SaintsView: view,
		Title:      "Saints",
		Today:      today(),
	})
}

func (server *Server) handle_stylesheet(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprint(response, stylesheet)
}

func (server *Server) render(response http.ResponseWriter, data PageData) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	bytes, err := server.render_template_bytes(data, "")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Write(bytes)
}

func (server *Server) render_day(response http.ResponseWriter, value string) {
	bytes, err := server.render_day_bytes(value, "")
	if err != nil {
		if errors.Is(err, errDateNotFound) {
			server.render_not_found(response, value)
		} else {
			server.render_error(response, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.Write(bytes)
}

func (server *Server) render_day_bytes(value string, root string) ([]byte, error) {
	view, found, err := db.DayViewByGregorianDate(server.conn, value)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errDateNotFound
	}

	hymns, err := db.HymnsByDayID(server.conn, view.Day.ID)
	if err != nil {
		return nil, err
	}

	return server.render_template_bytes(PageData{
		Active:       "day",
		DateValue:    value,
		DayView:      view,
		FastFree:     day_event_titles(view.FastFreeEvents, view.Day.FastFreePeriods),
		Fasts:        day_event_titles(view.FastEvents, view.Day.Fasts),
		Feasts:       day_event_titles(view.FeastEvents, view.Day.Feasts),
		FastingLevel: fasting_level(view.Day),
		HymnCount:    len(hymns),
		NextDate:     shift_date(value, 1),
		PrevDate:     shift_date(value, -1),
		Remembrances: day_event_titles(view.RemembranceEvents, view.Day.Remembrances),
		Title:        "Calendar",
		Today:        today(),
	}, root)
}

func (server *Server) render_error(response http.ResponseWriter, status int, message string) {
	response.WriteHeader(status)
	server.render(response, PageData{
		Active: "error",
		Error:  message,
		Title:  "Error",
		Today:  today(),
	})
}

func (server *Server) render_not_found(response http.ResponseWriter, value string) {
	response.WriteHeader(http.StatusNotFound)
	server.render(response, PageData{
		Active:    "not_found",
		DateValue: value,
		Error:     "No calendar data was found for " + value + ".",
		Title:     "Date Not Found",
		Today:     today(),
	})
}

func (server *Server) render_template_bytes(data PageData, root string) ([]byte, error) {
	data = server.page_defaults(data, root)

	buffer := bytes.Buffer{}
	if err := server.templates.ExecuteTemplate(&buffer, data.Active+"_page", data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (server *Server) page_defaults(data PageData, root string) PageData {
	if root == "" {
		data.HomeLink = "/"
		data.AssetPrefix = "/assets/"
		data.APIPrefix = "/api/date/"
		data.DayPrefix = "/date/"
		data.SaintsPrefix = "/saints/"
		data.ReadPrefix = "/readings/"
		data.HymnsPrefix = "/hymns/"
	} else {
		data.HomeLink = root + "index.html"
		data.AssetPrefix = root + "assets/"
		data.APIPrefix = root + "api/date/"
		data.APISuffix = ".json"
		data.DayPrefix = root + "dates/"
		data.SaintsPrefix = root + "saints/"
		data.ReadPrefix = root + "readings/"
		data.HymnsPrefix = root + "hymns/"
		data.LinkSuffix = "/index.html"
	}

	if data.Today == "" {
		data.Today = today()
	}

	return data
}

func route_date(request *http.Request, prefix string) (string, int, string, bool) {
	value := strings.TrimPrefix(request.URL.Path, prefix)
	if strings.Contains(value, "/") {
		return "", http.StatusNotFound, "not found", false
	}

	if _, err := time.Parse("2006-01-02", value); err != nil {
		return "", http.StatusBadRequest, "invalid date", false
	}

	return value, http.StatusOK, "", true
}

var errDateNotFound = errors.New("date not found")

func day_event_titles(events []db.CalendarDayEvent, fallback string) []string {
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

func fasting_level(day db.CalendarDay) string {
	if strings.TrimSpace(day.FastingLevelName) == "" {
		return ""
	}

	if strings.TrimSpace(day.FastingLevelCode) == "" {
		return day.FastingLevelName
	}

	return fmt.Sprintf("%s (%s)", day.FastingLevelName, day.FastingLevelCode)
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

func write_file(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}

func write_json_file(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(value)
}

func shift_date(value string, days int) string {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return value
	}

	return parsed.AddDate(0, 0, days).Format("2006-01-02")
}

func today() string {
	return time.Now().Format("2006-01-02")
}

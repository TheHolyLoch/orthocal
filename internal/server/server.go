// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/server.go

package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"orthocal/internal/db"
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
	mux := http.NewServeMux()
	server.routes(mux)

	return http.ListenAndServe(addr, mux)
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
	value, ok := route_date(response, request, "/date/")
	if !ok {
		return
	}

	server.render_day(response, value)
}

func (server *Server) handle_hymns(response http.ResponseWriter, request *http.Request) {
	value, ok := route_date(response, request, "/hymns/")
	if !ok {
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
	value, ok := route_date(response, request, "/readings/")
	if !ok {
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
	value, ok := route_date(response, request, "/saints/")
	if !ok {
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
	if err := server.templates.ExecuteTemplate(response, data.Active+"_page", data); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) render_day(response http.ResponseWriter, value string) {
	view, found, err := db.DayViewByGregorianDate(server.conn, value)
	if err != nil {
		server.render_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		server.render_not_found(response, value)
		return
	}

	hymns, err := db.HymnsByDayID(server.conn, view.Day.ID)
	if err != nil {
		server.render_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	server.render(response, PageData{
		Active:    "day",
		DateValue: value,
		DayView:   view,
		HymnCount: len(hymns),
		NextDate:  shift_date(value, 1),
		PrevDate:  shift_date(value, -1),
		Title:     "Calendar",
		Today:     today(),
	})
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

func route_date(response http.ResponseWriter, request *http.Request, prefix string) (string, bool) {
	value := strings.TrimPrefix(request.URL.Path, prefix)
	if strings.Contains(value, "/") {
		http.NotFound(response, request)
		return "", false
	}

	if _, err := time.Parse("2006-01-02", value); err != nil {
		http.Error(response, "invalid date", http.StatusBadRequest)
		return "", false
	}

	return value, true
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

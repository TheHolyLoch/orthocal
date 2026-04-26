// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/api.go

package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"orthocal/internal/db"
)

func (server *Server) handle_api_date(response http.ResponseWriter, request *http.Request) {
	value, ok := api_route_date(response, request, "/api/date/")
	if !ok {
		return
	}

	server.write_api_day(response, value)
}

func (server *Server) handle_api_hymns(response http.ResponseWriter, request *http.Request) {
	value, ok := api_route_date(response, request, "/api/hymns/")
	if !ok {
		return
	}

	view, found, err := db.HymnsViewByGregorianDate(server.conn, value)
	if err != nil {
		write_api_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		write_api_error(response, http.StatusNotFound, "date not found")
		return
	}

	write_json(response, http.StatusOK, view)
}

func (server *Server) handle_api_info(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/api/info" {
		http.NotFound(response, request)
		return
	}

	view, err := db.InfoViewByPath(server.conn, server.config.DatabasePath)
	if err != nil {
		write_api_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	write_json(response, http.StatusOK, view)
}

func (server *Server) handle_api_readings(response http.ResponseWriter, request *http.Request) {
	value, ok := api_route_date(response, request, "/api/readings/")
	if !ok {
		return
	}

	view, found, err := db.ReadingsViewByGregorianDate(server.conn, value)
	if err != nil {
		write_api_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		write_api_error(response, http.StatusNotFound, "date not found")
		return
	}

	write_json(response, http.StatusOK, view)
}

func (server *Server) handle_api_saints(response http.ResponseWriter, request *http.Request) {
	value, ok := api_route_date(response, request, "/api/saints/")
	if !ok {
		return
	}

	view, found, err := db.SaintsViewByGregorianDate(server.conn, value)
	if err != nil {
		write_api_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		write_api_error(response, http.StatusNotFound, "date not found")
		return
	}

	write_json(response, http.StatusOK, view)
}

func (server *Server) handle_api_today(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/api/today" {
		http.NotFound(response, request)
		return
	}

	server.write_api_day(response, today())
}

func (server *Server) handle_api_tomorrow(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/api/tomorrow" {
		http.NotFound(response, request)
		return
	}

	server.write_api_day(response, time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
}

func (server *Server) write_api_day(response http.ResponseWriter, value string) {
	view, found, err := db.DayViewByGregorianDate(server.conn, value)
	if err != nil {
		write_api_error(response, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		write_api_error(response, http.StatusNotFound, "date not found")
		return
	}

	write_json(response, http.StatusOK, view)
}

func api_route_date(response http.ResponseWriter, request *http.Request, prefix string) (string, bool) {
	value := strings.TrimPrefix(request.URL.Path, prefix)
	if strings.Contains(value, "/") {
		write_api_error(response, http.StatusNotFound, "not found")
		return "", false
	}

	if _, err := time.Parse("2006-01-02", value); err != nil {
		write_api_error(response, http.StatusBadRequest, "invalid date")
		return "", false
	}

	return value, true
}

func write_api_error(response http.ResponseWriter, status int, message string) {
	write_json(response, status, map[string]string{
		"error": message,
	})
}

func write_json(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	encoder := json.NewEncoder(response)
	encoder.SetIndent("", "\t")
	encoder.Encode(value)
}

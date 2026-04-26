// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/types.go

package server

import (
	"database/sql"
	"html/template"

	"orthocal/internal/db"
)

type Config struct {
	Addr         string
	DatabasePath string
}

type PageData struct {
	Active       string
	DateValue    string
	DayView      db.DayView
	Error        string
	HymnCount    int
	HymnsView    db.HymnsView
	NextDate     string
	PrevDate     string
	ReadingsView db.ReadingsView
	SaintsView   db.SaintsView
	Title        string
	Today        string
}

type Server struct {
	conn      *sql.DB
	config    Config
	templates *template.Template
}

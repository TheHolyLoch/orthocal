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
	APIPrefix    string
	APISuffix    string
	AppScript    template.JS
	AssetPrefix  string
	DateValue    string
	DayPrefix    string
	DayView      db.DayView
	Error        string
	HymnCount    int
	HymnsPrefix  string
	HymnsView    db.HymnsView
	HomeLink     string
	LinkSuffix   string
	NextDate     string
	PrevDate     string
	ReadingsView db.ReadingsView
	ReadPrefix   string
	SaintsView   db.SaintsView
	SaintsPrefix string
	StyleSheet   template.CSS
	Title        string
	Today        string
}

type Server struct {
	conn      *sql.DB
	config    Config
	templates *template.Template
}

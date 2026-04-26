// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/config/paths_test.go

package config

import (
	"path/filepath"
	"testing"
)

func TestResolveDBPathExplicitWins(t *testing.T) {
	t.Setenv("ORTHOCAL_DB", "/env/orthocal.db")
	t.Setenv("XDG_DATA_HOME", "/xdg")
	t.Setenv("HOME", "/home/test")

	got, err := ResolveDBPath("/explicit/orthocal.db")
	if err != nil {
		t.Fatalf("ResolveDBPath returned error: %v", err)
	}

	if got != "/explicit/orthocal.db" {
		t.Fatalf("expected explicit path, got %q", got)
	}
}

func TestResolveDBPathHomeFallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("ORTHOCAL_DB", "")
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("HOME", home)

	got, err := ResolveDBPath("")
	if err != nil {
		t.Fatalf("ResolveDBPath returned error: %v", err)
	}

	want := filepath.Join(home, ".local", "share", "orthocal", DefaultDatabaseName)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestResolveDBPathMissingHome(t *testing.T) {
	t.Setenv("ORTHOCAL_DB", "")
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("HOME", "")

	if _, err := ResolveDBPath(""); err == nil {
		t.Fatal("expected missing home error")
	}
}

func TestResolveDBPathOrthocalDB(t *testing.T) {
	t.Setenv("ORTHOCAL_DB", "/env/orthocal.db")
	t.Setenv("XDG_DATA_HOME", "/xdg")
	t.Setenv("HOME", "/home/test")

	got, err := ResolveDBPath("")
	if err != nil {
		t.Fatalf("ResolveDBPath returned error: %v", err)
	}

	if got != "/env/orthocal.db" {
		t.Fatalf("expected ORTHOCAL_DB path, got %q", got)
	}
}

func TestResolveDBPathXDGDataHome(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("ORTHOCAL_DB", "")
	t.Setenv("XDG_DATA_HOME", dataHome)
	t.Setenv("HOME", "/home/test")

	got, err := ResolveDBPath("")
	if err != nil {
		t.Fatalf("ResolveDBPath returned error: %v", err)
	}

	want := filepath.Join(dataHome, "orthocal", DefaultDatabaseName)
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

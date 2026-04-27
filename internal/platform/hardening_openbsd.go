//go:build openbsd

// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/platform/hardening_openbsd.go

package platform

import (
	"path/filepath"

	"golang.org/x/sys/unix"
)

func PledgeCLI() error {
	return unix.Pledge("stdio rpath", "")
}

func PledgeNetworkClient() error {
	return unix.Pledge("stdio rpath wpath cpath fattr inet dns", "")
}

func PledgeServer() error {
	return unix.Pledge("stdio rpath inet dns", "")
}

func UnveilNone() error {
	return unix.UnveilBlock()
}

func UnveilReadOnly(path string) error {
	return unveil_path(path, "r")
}

func UnveilReadWrite(path string) error {
	return unveil_path(path, "rwc")
}

func unveil_path(path string, flags string) error {
	if path == "" {
		return nil
	}

	if err := unix.Unveil(path, flags); err != nil {
		return err
	}

	absolute, err := filepath.Abs(path)
	if err != nil || absolute == path {
		return err
	}

	return unix.Unveil(absolute, flags)
}

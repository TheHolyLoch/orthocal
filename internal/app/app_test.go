// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/app/app_test.go

package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	code := Run([]string{"version"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%q", code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "program: orthocal") {
		t.Fatalf("unexpected version output: %q", stdout.String())
	}
}

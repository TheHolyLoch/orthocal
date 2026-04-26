// Orthocal - Developed by dgm (dgm@tuta.com)
// orthocal/internal/server/api_test.go

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIDate(t *testing.T) {
	server := test_server(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/date/2026-04-12", nil)

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	payload := map[string]any{}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json decode failed: %v", err)
	}

	day := payload["day"].(map[string]any)
	if day["gregorian_date"] != "2026-04-12" {
		t.Fatalf("unexpected gregorian date: %#v", day)
	}

	primarySaints := payload["primary_saints"].([]any)
	if len(primarySaints) == 0 {
		t.Fatal("expected at least one primary saint")
	}
}

func TestAPIInfo(t *testing.T) {
	server := test_server(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/info", nil)

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if !strings.Contains(response.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("expected JSON content type, got %q", response.Header().Get("Content-Type"))
	}

	payload := map[string]any{}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json decode failed: %v", err)
	}

	counts := payload["counts"].(map[string]any)
	if counts["calendar_days"].(float64) != 1 {
		t.Fatalf("unexpected counts: %#v", counts)
	}
}

func TestAPIInvalidDate(t *testing.T) {
	server := test_server(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/date/not-a-date", nil)

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), `"error"`) {
		t.Fatalf("expected error JSON, got %s", response.Body.String())
	}
}

func TestAPIMissingDate(t *testing.T) {
	server := test_server(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/date/2099-01-01", nil)

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), `"error"`) {
		t.Fatalf("expected error JSON, got %s", response.Body.String())
	}
}

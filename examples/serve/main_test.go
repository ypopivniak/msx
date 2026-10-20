package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MSX setup loads {SERVER}/msx/start.json, and every other document must be
// reachable through the absolute URLs it links to.
func TestHandler_ServesSetupLayout(t *testing.T) {
	base := "http://192.168.0.10:8080/msx"
	h, err := newHandler(base)
	if err != nil {
		t.Fatalf("newHandler: %v", err)
	}

	for path, want := range map[string]string{
		"/msx/start.json":       `"parameter": "menu:` + base + `/data/menu.json"`,
		"/msx/data/menu.json":   `"data": "` + base + `/data/videos.json"`,
		"/msx/data/videos.json": `"headline": "Videos"`,
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		body, _ := io.ReadAll(rec.Body)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status %d", path, rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
			t.Errorf("%s: Access-Control-Allow-Origin = %q", path, got)
		}
		if !strings.Contains(string(body), want) {
			t.Errorf("%s: expected %s in\n%s", path, want, body)
		}
	}
}

func TestServerAddress(t *testing.T) {
	for _, tc := range []struct{ addr, host, want string }{
		{":8080", "tv.lan", "tv.lan:8080"},
		{"192.168.0.10:8080", "", "192.168.0.10:8080"},
		{"0.0.0.0:9000", "example.com", "example.com:9000"},
	} {
		got, err := serverAddress(tc.addr, tc.host)
		if err != nil || got != tc.want {
			t.Errorf("serverAddress(%q, %q) = %q, %v; want %q", tc.addr, tc.host, got, err, tc.want)
		}
	}

	// Unspecified hosts fall back to a detected LAN IP, never "localhost"
	// or the bare wildcard address.
	for _, addr := range []string{":8080", "0.0.0.0:8080"} {
		got, err := serverAddress(addr, "")
		if err != nil {
			t.Skipf("no LAN address on this machine: %v", err)
		}
		if strings.HasPrefix(got, "0.0.0.0") || strings.HasPrefix(got, ":") || strings.HasPrefix(got, "localhost") {
			t.Errorf("serverAddress(%q) = %q, want a LAN host", addr, got)
		}
	}

	if _, err := serverAddress("8080", ""); err == nil {
		t.Error("expected error for addr without port separator")
	}
}

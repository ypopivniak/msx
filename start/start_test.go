package start_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/start"
)

// TestStart_BuilderReproducesReference rebuilds the authoritative start.json
// from benzac-de/msx using the fluent API.
func TestStart_BuilderReproducesReference(t *testing.T) {
	built, err := start.New().
		WithName("GitHub Pages Example").
		WithVersion("1.0.3").
		WithParameter("menu:https://benzac-de.github.io/msx/data/menu.json").
		WithLauncherIcon("pages").
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	want := mustRead(t, "../testdata/reference/start.json")
	got, err := json.Marshal(built)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	assertJSONSemanticEq(t, got, want)
}

// TestStart_RoundTrip unmarshals the reference start.json, re-marshals via
// our types, and compares.
func TestStart_RoundTrip(t *testing.T) {
	raw := mustRead(t, "../testdata/reference/start.json")

	var s start.Start
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	got, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	assertJSONSemanticEq(t, got, raw)
}

func TestStart_RejectsMissingRequired(t *testing.T) {
	cases := map[string]*start.Start{
		"name":      start.New().WithVersion("1").WithParameter("menu:x"),
		"version":   start.New().WithName("x").WithParameter("menu:x"),
		"parameter": start.New().WithName("x").WithVersion("1"),
	}
	for field, s := range cases {
		field, s := field, s
		t.Run(field, func(t *testing.T) {
			_, err := s.Build()
			if err == nil {
				t.Fatalf("expected error for missing %s", field)
			}
			if !strings.Contains(err.Error(), field+": is required") {
				t.Errorf("expected 'is required' on %s, got: %v", field, err)
			}
		})
	}
}

func TestStart_RejectsBadWelcome(t *testing.T) {
	_, err := start.New().
		WithName("x").WithVersion("1").WithParameter("menu:x").
		WithWelcome(start.Welcome("greeting")).
		Build()
	if err == nil {
		t.Fatal("expected welcome enum error, got nil")
	}
	if !strings.Contains(err.Error(), "welcome") {
		t.Errorf("expected welcome in error, got: %v", err)
	}
}

func TestStart_LauncherBuilderIndependent(t *testing.T) {
	l := start.NewLauncher().
		WithIcon(icon.Home).
		WithColor("msx-red").
		WithImage("https://example/tile.png")
	built, err := start.New().
		WithName("x").WithVersion("1").WithParameter("menu:x").
		WithLauncher(l).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	b, _ := json.Marshal(built)
	if !bytes.Contains(b, []byte(`"icon":"home"`)) {
		t.Errorf("expected icon=home in output, got %s", b)
	}
	if !bytes.Contains(b, []byte(`"color":"msx-red"`)) {
		t.Errorf("expected color=msx-red in output, got %s", b)
	}
}

// ---- helpers ----

func mustRead(t *testing.T, rel string) []byte {
	t.Helper()
	p, err := filepath.Abs(rel)
	if err != nil {
		t.Fatalf("abs %s: %v", rel, err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", p, err)
	}
	return data
}

func assertJSONSemanticEq(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	gc, _ := json.Marshal(g)
	wc, _ := json.Marshal(w)
	if !bytes.Equal(gc, wc) {
		t.Errorf("JSON mismatch:\n got:  %s\n want: %s", gc, wc)
	}
}

package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// Each builder must succeed and produce a non-empty JSON document with
// the expected top-level shape. This is an integration smoke test that
// catches cross-package breakage without requiring the example binary
// to be executed.
func TestBuilders_Integrate(t *testing.T) {
	base := "https://example.invalid/msx"

	s, err := buildStart(base)
	if err != nil {
		t.Fatalf("buildStart: %v", err)
	}
	if b, _ := json.Marshal(s); !strings.Contains(string(b), `"parameter":"menu:`) {
		t.Errorf("start parameter missing: %s", b)
	}

	m, err := buildMenu(base)
	if err != nil {
		t.Fatalf("buildMenu: %v", err)
	}
	if b, _ := json.Marshal(m); !strings.Contains(string(b), `"menu":[`) {
		t.Errorf("menu shape wrong: %s", b)
	} else if !strings.Contains(string(b), `"dictionary":"`+base+`/data/dictionary.json"`) {
		t.Errorf("menu should reference the generated dictionary: %s", b)
	}

	c, err := buildVideos()
	if err != nil {
		t.Fatalf("buildVideos: %v", err)
	}
	b, _ := json.Marshal(c)
	if !strings.Contains(string(b), `"{ico:msx-white:movie} Videos"`) {
		t.Errorf("expected expr-built headline in videos content: %s", b)
	}
	if !strings.Contains(string(b), `"type":"separate"`) {
		t.Errorf("expected separate template type: %s", b)
	}

	d, err := buildDictionary()
	if err != nil {
		t.Fatalf("buildDictionary: %v", err)
	}
	if b, _ := json.Marshal(d); !strings.Contains(string(b), `"dictionary:message":[`) {
		t.Errorf("expected array-valued dictionary entry: %s", b)
	}
}

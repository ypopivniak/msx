package dictionary_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/dictionary"
)

func TestDictionary_SingleLineShape(t *testing.T) {
	d, err := dictionary.New().
		WithName("English").
		WithVersion("1.0.35").
		WithLanguage("en").
		Set("caption:menu", "menu").
		Set("label:about", "About").
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"name":"English"`,
		`"version":"1.0.35"`,
		`"language":"en"`,
		`"caption:menu":"menu"`,
		`"label:about":"About"`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in %s", want, s)
		}
	}
}

func TestDictionary_MultilineValue(t *testing.T) {
	d, err := dictionary.New().
		WithName("English").
		WithVersion("1").
		SetLines("message:welcome", "Hello.", "Second line.").
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	b, _ := json.Marshal(d)
	if !strings.Contains(string(b), `"message:welcome":["Hello.","Second line."]`) {
		t.Errorf("expected array-valued entry, got %s", b)
	}
}

func TestDictionary_RejectsMissingRequired(t *testing.T) {
	cases := map[string]*dictionary.Dictionary{
		"name":       dictionary.New().WithVersion("1").Set("k", "v"),
		"version":    dictionary.New().WithName("x").Set("k", "v"),
		"properties": dictionary.New().WithName("x").WithVersion("1"),
	}
	for field, d := range cases {
		field, d := field, d
		t.Run(field, func(t *testing.T) {
			_, err := d.Build()
			if err == nil {
				t.Fatalf("expected error for missing %s", field)
			}
			if !strings.Contains(err.Error(), field) {
				t.Errorf("expected %s in error, got %v", field, err)
			}
		})
	}
}

// Round-trip: author a dictionary, marshal, unmarshal, verify values and
// that single-vs-array shapes are preserved.
func TestDictionary_RoundTrip(t *testing.T) {
	orig, err := dictionary.New().
		WithName("English").WithVersion("1").WithLanguage("en").
		Set("caption:menu", "menu").
		SetLines("message:about", "Line A", "Line B").
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	raw, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var back dictionary.Dictionary
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v := back.Properties["caption:menu"]; v.String() != "menu" || v.Lines() != nil {
		t.Errorf("single-line roundtrip failed: %#v", v)
	}
	if v := back.Properties["message:about"]; v.String() != "" ||
		len(v.Lines()) != 2 || v.Lines()[0] != "Line A" || v.Lines()[1] != "Line B" {
		t.Errorf("multi-line roundtrip failed: %#v", v)
	}
}

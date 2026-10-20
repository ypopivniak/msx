package menu_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/menu"
)

// TestMenu_RoundTrip unmarshals the real benzac-de/msx menu.json, re-marshals
// via our types, and asserts the result is semantically equal.
func TestMenu_RoundTrip(t *testing.T) {
	raw := mustReadBytes(t, "../testdata/reference/data/menu.json")

	var m menu.Menu
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	got, err := json.Marshal(&m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	assertJSONSemanticEq(t, got, raw)
}

// TestMenu_BuilderReproduces_Reference rebuilds the reference menu.json
// with the fluent API, exercising separators, ref-style data, and icons.
func TestMenu_BuilderReproduces_Reference(t *testing.T) {
	ref := func(label string, ic icon.Icon, url string) *menu.Item {
		return menu.NewItem().
			WithIcon(ic).
			WithLabel(label).
			WithDataRef(url)
	}
	sep := func(label string) *menu.Item {
		return menu.NewItem().WithType(menu.ItemTypeSeparator).WithLabel(label)
	}

	base := "https://benzac-de.github.io/msx"

	built, err := menu.New().
		WithHeadline("GitHub Pages Example").
		AddItem(ref("Welcome", icon.Info, base+"/data/welcome.json")).
		AddItem(sep("Base")).
		AddItem(ref("Videos", icon.Movie, base+"/data/base/videos.json")).
		AddItem(ref("Audios", icon.MusicNote, base+"/data/base/audios.json")).
		AddItem(ref("Images", icon.Image, base+"/data/base/images.json")).
		AddItem(ref("YouTube", "smart-display", base+"/data/base/youtube.json")).
		AddItem(sep("Extended - Lists")).
		AddItem(ref("First List", "view-stream", base+"/data/extended/list0.json")).
		AddItem(ref("Second List", "view-stream", base+"/data/extended/list1.json")).
		AddItem(ref("Third List", "view-stream", base+"/data/extended/list2.json")).
		AddItem(ref("Fourth List", "view-stream", base+"/data/extended/list3_section0.json")).
		AddItem(ref("Fifth List", "view-stream", base+"/data/extended/list4.json")).
		AddItem(ref("Sixth List", "view-stream", base+"/data/extended/list5.json")).
		AddItem(sep("Extended - Pages")).
		AddItem(ref("First Pages", "view-array", base+"/data/extended/pages0.json")).
		AddItem(ref("Second Pages", "view-array", base+"/data/extended/pages1.json")).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	raw := mustReadBytes(t, "../testdata/reference/data/menu.json")
	got, err := json.Marshal(built)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	assertJSONSemanticEq(t, got, raw)
}

// ---- helpers (local to this file to avoid reliance on the other test file) ----

func mustReadBytes(t *testing.T, rel string) []byte {
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
		t.Fatalf("unmarshal got: %v\n%s", err, got)
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatalf("unmarshal want: %v\n%s", err, want)
	}
	gc, _ := json.Marshal(g)
	wc, _ := json.Marshal(w)
	if !bytes.Equal(gc, wc) {
		t.Errorf("JSON mismatch:\n got:  %s\n want: %s", gc, wc)
	}
}

package content_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/content"
)

// contentFixtures enumerates the reference JSON documents from benzac-de/msx
// that are Content Root Objects (not Menu Root Objects). Each one must
// round-trip through our types without data loss.
var contentFixtures = []string{
	"../testdata/reference/data/welcome.json",
	"../testdata/reference/data/base/videos.json",
	"../testdata/reference/data/base/audios.json",
	"../testdata/reference/data/base/images.json",
	"../testdata/reference/data/base/youtube.json",
	"../testdata/reference/data/extended/pages0.json",
	"../testdata/reference/data/extended/pages1.json",
}

// TestContent_RoundTrip unmarshals every content fixture into content.Content,
// re-marshals it, and asserts the result is semantically equal to the source.
// Any property the library doesn't model would either fail to decode or get
// dropped on re-encode, and the canonical compare would fail.
func TestContent_RoundTrip(t *testing.T) {
	for _, rel := range contentFixtures {
		rel := rel
		t.Run(filepath.Base(rel), func(t *testing.T) {
			raw := mustRead(t, rel)

			var c content.Content
			if err := json.Unmarshal(raw, &c); err != nil {
				t.Fatalf("unmarshal %s: %v", rel, err)
			}

			// After a round-trip the result must also pass strict validation.
			if err := c.Validate(); err != nil {
				t.Fatalf("validate %s: %v", rel, err)
			}

			got, err := json.Marshal(&c)
			if err != nil {
				t.Fatalf("marshal %s: %v", rel, err)
			}

			assertJSONSemanticallyEqual(t, got, raw)
		})
	}
}

// TestContent_BuilderReproduces_Videos exercises the fluent API end-to-end
// for videos.json — a compact Content Root with a template and three items.
func TestContent_BuilderReproduces_Videos(t *testing.T) {
	template := content.NewItem().
		WithType(content.ItemTypeSeparate).
		WithLayoutXYWH(0, 0, 2, 4).
		WithIcon("msx-white-soft:movie").
		WithColor("msx-glass")

	video := func(n int, url string) *content.Item {
		title := "Video " + strconv.Itoa(n)
		return content.NewItem().
			WithTitle(title).
			WithPlayerLabel(title).
			WithAction(action.Video(url))
	}

	built, err := content.New().
		WithType(content.TypePages).
		WithHeadline("Videos").
		WithTemplate(template).
		WithItems(
			video(1, "http://msx.benzac.de/media/video1.mp4"),
			video(2, "http://msx.benzac.de/media/video2.mp4"),
			video(3, "http://msx.benzac.de/media/video3.mp4"),
		).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	raw := mustRead(t, "../testdata/reference/data/base/videos.json")
	got, err := json.Marshal(built)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	assertJSONSemanticallyEqual(t, got, raw)
}

// TestContent_BuilderReproduces_Welcome exercises the text (array form),
// default item type, and inline-expression-bearing strings.
func TestContent_BuilderReproduces_Welcome(t *testing.T) {
	built, err := content.New().
		AddPage(content.NewPage().AddItem(
			content.NewItem().
				WithType(content.ItemTypeDefault).
				WithLayoutXYWH(0, 0, 12, 6).
				WithColor("msx-glass").
				WithHeadline("Welcome to Media Station X {chr:ndash} GitHub Pages Example").
				WithText(
					"This project shows a way to host {txt:msx-white:Media Station X} JSON files via {txt:msx-white:GitHub Pages}.{br}{br}",
					"{ico:msx-white:play-arrow} For more information about the {txt:msx-white:Media Station X} project, please visit: {txt:msx-white:https://msx.benzac.de/info/}.{br}",
					"{ico:msx-white:play-arrow} Check out this project on GitHub: {txt:msx-white:https://github.com/benzac-de/msx}.{br}",
				),
		)).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	raw := mustRead(t, "../testdata/reference/data/welcome.json")
	got, err := json.Marshal(built)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	assertJSONSemanticallyEqual(t, got, raw)
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

// assertJSONSemanticallyEqual compares two JSON payloads ignoring formatting,
// whitespace, and key order. It also strips keys whose value is an empty
// string at any depth: MSX treats "" as the default for optional string
// properties, so fixtures that explicitly write "headline": "" round-trip
// through our omitempty output as a missing key, which is semantically the
// same thing.
func assertJSONSemanticallyEqual(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("unmarshal got: %v\n%s", err, got)
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatalf("unmarshal want: %v\n%s", err, want)
	}
	g = stripEmptyStrings(g)
	w = stripEmptyStrings(w)
	gc, _ := json.Marshal(g)
	wc, _ := json.Marshal(w)
	if !bytes.Equal(gc, wc) {
		t.Errorf("JSON mismatch:\n got:  %s\n want: %s", gc, wc)
	}
}

// stripEmptyStrings walks a decoded JSON tree and removes keys whose value
// is an empty string, normalizing "explicit empty defaults" that MSX
// fixtures sometimes emit but which omitempty rightfully drops on output.
func stripEmptyStrings(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, vv := range x {
			if s, ok := vv.(string); ok && s == "" {
				continue
			}
			out[k] = stripEmptyStrings(vv)
		}
		return out
	case []any:
		for i, vv := range x {
			x[i] = stripEmptyStrings(vv)
		}
		return x
	default:
		return v
	}
}


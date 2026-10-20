package content_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/content"
)

// A minimal "pages with one teaser" Content Root builds to the expected JSON.
func TestBuilder_ReproducesMinimalPages(t *testing.T) {
	got, err := content.New().
		WithType(content.TypePages).
		AddPage(content.NewPage().AddItem(
			content.NewItem().
				WithType(content.ItemTypeTeaser).
				WithLayoutXYWH(0, 0, 4, 6).
				WithTitle("Video").
				WithAction(action.Video("https://example.invalid/index.m3u8")),
		)).
		Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	want := `{
      "type": "pages",
      "pages": [
        {
          "items": [
            {
              "type": "teaser",
              "layout": "0,0,4,6",
              "title": "Video",
              "action": "video:https://example.invalid/index.m3u8"
            }
          ]
        }
      ]
    }`

	assertJSONEqual(t, got, want)
}

func TestBuilder_RejectsTemplateAndPagesTogether(t *testing.T) {
	_, err := content.New().
		WithTemplate(content.NewItem().WithType(content.ItemTypeTeaser)).
		AddPage(content.NewPage().AddItem(content.NewItem())).
		Build()
	if err == nil {
		t.Fatal("expected error for template+pages together, got nil")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("expected mutual-exclusion error, got: %v", err)
	}
}

func TestBuilder_RejectsTemplateWithoutItems(t *testing.T) {
	_, err := content.New().
		WithTemplate(content.NewItem().WithType(content.ItemTypeTeaser)).
		Build()
	if err == nil {
		t.Fatal("expected error for template without items, got nil")
	}
	if !strings.Contains(err.Error(), "required when template is set") {
		t.Errorf("expected required-items error, got: %v", err)
	}
}

func TestBuilder_RejectsEmptyPage(t *testing.T) {
	_, err := content.New().AddPage(content.NewPage()).Build()
	if err == nil {
		t.Fatal("expected error for empty page, got nil")
	}
	if !strings.Contains(err.Error(), "at least one item") {
		t.Errorf("expected empty-items error, got: %v", err)
	}
}

func TestItem_RejectsBadLayout(t *testing.T) {
	_, err := content.New().
		AddPage(content.NewPage().AddItem(
			content.NewItem().WithLayoutXYWH(99, 0, 4, 6),
		)).
		Build()
	if err == nil {
		t.Fatal("expected layout.x range error, got nil")
	}
	if !strings.Contains(err.Error(), "layout.x") {
		t.Errorf("expected layout.x error path, got: %v", err)
	}
}

func TestItem_RejectsBadEnum(t *testing.T) {
	_, err := content.New().
		AddPage(content.NewPage().AddItem(
			&content.Item{Type: content.ItemType("banana")},
		)).
		Build()
	if err == nil {
		t.Fatal("expected enum error, got nil")
	}
	if !strings.Contains(err.Error(), "items[0].type") {
		t.Errorf("expected items[0].type path, got: %v", err)
	}
}

func TestItem_RejectsBadID(t *testing.T) {
	_, err := content.New().
		AddPage(content.NewPage().AddItem(
			content.NewItem().WithID("has spaces"),
		)).
		Build()
	if err == nil {
		t.Fatal("expected id format error, got nil")
	}
}

func TestItem_ExplicitFalseForDefaultTrueBool(t *testing.T) {
	got, err := content.New().
		AddPage(content.NewPage().AddItem(
			content.NewItem().WithDisplay(false).WithTitle("T"),
		)).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// display=false must survive the round-trip even though bool zero is false.
	assertJSONContains(t, got, `"display":false`)
}

func TestItem_TextSingleLineMarshalsAsString(t *testing.T) {
	it := content.NewItem().WithText("hello")
	b, err := json.Marshal(it.Text)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != `"hello"` {
		t.Errorf("single-line text should marshal as string, got %s", b)
	}
}

func TestItem_TextMultiLineMarshalsAsArray(t *testing.T) {
	it := content.NewItem().WithText("a", "b")
	b, err := json.Marshal(it.Text)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(b) != `["a","b"]` {
		t.Errorf("multi-line text should marshal as array, got %s", b)
	}
}

// ---- helpers ----

// assertJSONEqual compares two JSON documents for semantic equality, ignoring
// whitespace and key ordering.
func assertJSONEqual(t *testing.T, got any, wantJSON string) {
	t.Helper()
	gotBytes, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	var gotV, wantV any
	if err := json.Unmarshal(gotBytes, &gotV); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	if err := json.Unmarshal([]byte(wantJSON), &wantV); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	gotCanon, _ := json.Marshal(gotV)
	wantCanon, _ := json.Marshal(wantV)
	if !bytes.Equal(gotCanon, wantCanon) {
		t.Errorf("JSON mismatch:\n got:  %s\n want: %s", gotCanon, wantCanon)
	}
}

func assertJSONContains(t *testing.T, v any, fragment string) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), fragment) {
		t.Errorf("expected JSON to contain %q, got: %s", fragment, b)
	}
}

// refocus/transparent are documented as number|boolean; the boolean form
// must decode (false → 0, true → 1) and re-emit as a number.
func TestTri_AcceptsBooleanForm(t *testing.T) {
	var c content.Content
	raw := `{"refocus":true,"transparent":2,"pages":[{"transparent":true,"items":[{"layout":"0,0,1,1"}]}]}`
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Refocus != content.TriOn || c.Transparent != content.TriAuto || c.Pages[0].Transparent != content.TriOn {
		t.Fatalf("got refocus=%d transparent=%d page.transparent=%d", c.Refocus, c.Transparent, c.Pages[0].Transparent)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	assertJSONEqual(t, &c, `{"refocus":1,"transparent":2,"pages":[{"transparent":1,"items":[{"layout":"0,0,1,1"}]}]}`)

	var off content.Content
	if err := json.Unmarshal([]byte(`{"refocus":false}`), &off); err != nil {
		t.Fatalf("unmarshal false: %v", err)
	}
	if off.Refocus != content.TriOff {
		t.Errorf("refocus=false decoded as %d, want 0", off.Refocus)
	}
	if err := json.Unmarshal([]byte(`{"refocus":"yes"}`), &off); err == nil {
		t.Error("expected error for string refocus")
	}
}

// Extended properties accept string, number, and boolean values and keep
// their JSON type through a round-trip.
func TestItem_PropertiesNonStringValues(t *testing.T) {
	raw := `{"layout":"0,0,1,1","properties":{"control:dim":true,"resume:position":123,"button:content:icon":"info"}}`
	var it content.Item
	if err := json.Unmarshal([]byte(raw), &it); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := it.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	assertJSONEqual(t, &it, raw)

	built := content.NewItem().
		WithProperty("button:content:icon", "info").
		WithPropertyValue("control:dim", true)
	assertJSONEqual(t, built, `{"properties":{"button:content:icon":"info","control:dim":true}}`)

	legacy := content.NewItem().WithProperties(map[string]string{"control:dim": "true"})
	assertJSONEqual(t, legacy, `{"properties":{"control:dim":"true"}}`)

	bad := content.NewItem().WithPropertyValue("x", map[string]any{"nested": 1})
	if err := bad.Validate(); err == nil || !strings.Contains(err.Error(), "properties.x") {
		t.Errorf("expected properties.x error, got %v", err)
	}
}

// -1 is the documented default for separation/imageWidth/imageHeight and
// any negative progress hides the bar, so explicit values must validate;
// genuinely out-of-range values are still rejected.
func TestItem_AcceptsDocumentedSentinels(t *testing.T) {
	ok := content.NewItem().WithLayoutXYWH(0, 0, 1, 1).
		WithSeparation(-1).WithImageWidth(-1).WithImageHeight(-1).WithProgress(-1)
	if err := ok.Validate(); err != nil {
		t.Fatalf("Validate sentinels: %v", err)
	}
	if err := content.NewItem().WithProgress(-0.5).Validate(); err != nil {
		t.Errorf("negative progress: %v", err)
	}

	for name, it := range map[string]*content.Item{
		"separation":  content.NewItem().WithSeparation(-0.5),
		"imageWidth":  content.NewItem().WithImageWidth(13),
		"imageHeight": content.NewItem().WithImageHeight(-2),
		"progress":    content.NewItem().WithProgress(1.5),
	} {
		err := it.Validate()
		if err == nil || !strings.Contains(err.Error(), name) {
			t.Errorf("%s: expected error, got %v", name, err)
		}
	}
}

// An options Content Root that has no pages/name/version (e.g. a templated
// list) must still decode as a Root, so root-only keys like `type` survive.
func TestPanel_DetectsRootWithoutPages(t *testing.T) {
	raw := `{"options":{"type":"list","headline":"H","template":{"layout":"0,0,12,1"},"items":[{"title":"a"}]}}`
	var it content.Item
	if err := json.Unmarshal([]byte(raw), &it); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if it.Options.Root() == nil {
		t.Fatalf("options decoded as page, want root")
	}
	if err := it.Options.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	assertJSONEqual(t, &it, raw)

	var page content.Item
	if err := json.Unmarshal([]byte(`{"options":{"headline":"H","items":[{"layout":"0,0,1,1"}]}}`), &page); err != nil {
		t.Fatalf("unmarshal page: %v", err)
	}
	if page.Options.Page() == nil {
		t.Error("plain page options decoded as root")
	}
}

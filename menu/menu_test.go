package menu_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/menu"
)

func TestBuilder_RequiresNonEmptyMenu(t *testing.T) {
	_, err := menu.New().WithName("X").Build()
	if err == nil {
		t.Fatal("expected error for empty menu, got nil")
	}
	if !strings.Contains(err.Error(), "at least one item") {
		t.Errorf("expected at-least-one-item error, got: %v", err)
	}
}

func TestBuilder_RejectsBadStyle(t *testing.T) {
	_, err := menu.New().
		WithStyle(menu.Style("fancy")).
		AddItem(menu.NewItem().WithLabel("X")).
		Build()
	if err == nil {
		t.Fatal("expected style enum error, got nil")
	}
	if !strings.Contains(err.Error(), "style") {
		t.Errorf("expected style path in error, got: %v", err)
	}
}

func TestItem_RejectsLineColorOnNonSeparator(t *testing.T) {
	_, err := menu.New().
		AddItem(menu.NewItem().
			WithLabel("X").
			WithLineColor(color.MSXRed)).
		Build()
	if err == nil {
		t.Fatal("expected lineColor error, got nil")
	}
	if !strings.Contains(err.Error(), "lineColor") {
		t.Errorf("expected lineColor in error, got: %v", err)
	}
}

func TestItemData_RejectsBothRefAndInline(t *testing.T) {
	d := &menu.ItemData{
		Ref:    "content:http://example/",
		Inline: &content.Content{},
	}
	_, err := menu.New().
		AddItem(&menu.Item{Label: "X", Data: d}).
		Build()
	if err == nil {
		t.Fatal("expected mutual-exclusion error, got nil")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("expected exclusivity error, got: %v", err)
	}
}

func TestItemData_RefMarshalsAsString(t *testing.T) {
	d := menu.NewItemDataRef("content:http://example/")
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != `"content:http://example/"` {
		t.Errorf("expected ref to marshal as string, got %s", b)
	}
}

// ---- helpers ----

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

// Menu root and menu item refocus/transparent accept the boolean form.
func TestMenu_TriAcceptsBooleanForm(t *testing.T) {
	raw := `{"refocus":true,"transparent":false,"menu":[{"label":"A","transparent":true,"data":"content:https://example.com/a.json"}]}`
	var m menu.Menu
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m.Refocus != menu.TriOn || m.Transparent != menu.TriOff || m.Menu[0].Transparent != menu.TriOn {
		t.Fatalf("got refocus=%d transparent=%d item.transparent=%d", m.Refocus, m.Transparent, m.Menu[0].Transparent)
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	assertJSONEqual(t, &m, `{"refocus":1,"menu":[{"label":"A","transparent":1,"data":"content:https://example.com/a.json"}]}`)
}

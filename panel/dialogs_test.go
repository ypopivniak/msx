package panel_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/panel"
)

func TestInfo_PanelShape(t *testing.T) {
	p := panel.Info("Heads up", "This is informational.")
	if err := p.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	// Page-backed panel: no "pages" wrapper.
	if strings.Contains(s, `"pages"`) {
		t.Errorf("Info should marshal as a page, not a root: %s", s)
	}
	// Two items: one message block + one OK button.
	if !strings.Contains(s, `"headline":"Heads up"`) {
		t.Errorf("headline missing: %s", s)
	}
	if !strings.Contains(s, `"text":"This is informational."`) {
		t.Errorf("text missing: %s", s)
	}
	if !strings.Contains(s, `"label":"OK"`) {
		t.Errorf("OK button missing: %s", s)
	}
	if !strings.Contains(s, `"icon":"info"`) {
		t.Errorf("expected info icon for Info tone: %s", s)
	}
	if !strings.Contains(s, `"color":"msx-blue"`) {
		t.Errorf("expected msx-blue color for Info tone: %s", s)
	}
}

func TestError_UsesErrorTone(t *testing.T) {
	p := panel.Error("Oops", "Something failed.")
	b, _ := json.Marshal(p)
	s := string(b)
	if !strings.Contains(s, `"icon":"error"`) || !strings.Contains(s, `"color":"msx-red"`) {
		t.Errorf("expected error tone, got %s", s)
	}
}

func TestConfirm_HasYesAndNoButtons(t *testing.T) {
	onYes := action.Compound(
		action.ReloadScope("content"),
		action.Info("confirmed"),
	)
	p := panel.Confirm("Delete?", "This cannot be undone.", onYes, action.Back)
	if err := p.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	b, _ := json.Marshal(p)
	s := string(b)
	if !strings.Contains(s, `"label":"Yes"`) {
		t.Errorf("missing Yes button: %s", s)
	}
	if !strings.Contains(s, `"label":"No"`) {
		t.Errorf("missing No button: %s", s)
	}
	if !strings.Contains(s, `"focus":true`) {
		t.Errorf("expected Yes to have focus, got: %s", s)
	}
	if !strings.Contains(s, string(onYes)) {
		t.Errorf("expected compound confirm action in output: %s", s)
	}
}

func TestConfirm_DefaultsOnCancelToBack(t *testing.T) {
	p := panel.Confirm("X", "Y", action.Reload, "")
	b, _ := json.Marshal(p)
	s := string(b)
	// The No button should carry action "back".
	if !strings.Contains(s, `"label":"No","action":"back"`) &&
		!strings.Contains(s, `"action":"back","label":"No"`) {
		t.Errorf("No button should default to action:back, got %s", s)
	}
}

func TestMessage_CustomTone(t *testing.T) {
	tone := panel.MessageOption{Icon: "favorite", Color: "msx-pink"}
	p := panel.Message("Hi", "Custom", tone)
	b, _ := json.Marshal(p)
	s := string(b)
	if !strings.Contains(s, `"icon":"favorite"`) || !strings.Contains(s, `"color":"msx-pink"`) {
		t.Errorf("custom tone not applied: %s", s)
	}
}

package content_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/content"
)

// A schedule Live with coming/running/over states marshals with the
// expected field names and enum values.
func TestLive_ScheduleShape(t *testing.T) {
	live := content.NewLive().
		WithType(content.LiveSchedule).
		WithFrom(1_700_000_000_000).
		WithTo(1_700_003_600_000).
		WithComing(content.NewLiveState().
			WithHeadline("Starting soon").
			WithAction(action.Info("coming up"))).
		WithRunning(content.NewLiveRunningState().
			WithAction(action.Reload).
			WithQuartile2(content.NewLiveState().
				WithHeadline("Halfway"))).
		WithOver(content.NewLiveState().
			WithHeadline("Done").
			WithAction(action.Back))

	if err := live.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}

	b, err := json.Marshal(live)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	checks := []string{
		`"type":"schedule"`,
		`"from":1700000000000`,
		`"to":1700003600000`,
		`"coming":{`,
		`"running":{`,
		`"quartile2":{`,
		`"over":{`,
		`"action":"back"`,
	}
	for _, want := range checks {
		if !strings.Contains(s, want) {
			t.Errorf("expected %q in output, got %s", want, s)
		}
	}
}

func TestLive_RejectsBadType(t *testing.T) {
	live := content.NewLive().WithType(content.LiveType("nope"))
	if err := live.Validate(); err == nil {
		t.Fatal("expected enum error, got nil")
	}
}

func TestLive_RejectsNegativeTiming(t *testing.T) {
	live := &content.Live{From: -1}
	err := live.Validate()
	if err == nil || !strings.Contains(err.Error(), "from") {
		t.Errorf("expected from:>=0 error, got %v", err)
	}
}

func TestLive_RejectsToBeforeFrom(t *testing.T) {
	live := content.NewLive().WithFrom(100).WithTo(50)
	err := live.Validate()
	if err == nil || !strings.Contains(err.Error(), "to") {
		t.Errorf("expected to>=from error, got %v", err)
	}
}

// A Live attached to an Item recurses through Item.Validate.
func TestItem_WithLive_ValidatesNested(t *testing.T) {
	bad := content.NewLive().WithType(content.LiveType("unknown"))
	item := content.NewItem().WithTitle("x").WithLive(bad)
	err := item.Validate()
	if err == nil || !strings.Contains(err.Error(), "live.type") {
		t.Errorf("expected live.type path in error, got %v", err)
	}
}

// Selection typed: marshal path and nested in Item.
func TestSelection_MarshalShape(t *testing.T) {
	sel := content.NewSelection().
		WithImportant(true).
		WithHeadline("Focus me").
		WithAction(action.Compound(
			action.Update("content", "underlay:image"),
			action.Update("content", "underlay:text"),
		))
	b, _ := json.Marshal(sel)
	s := string(b)
	if !strings.Contains(s, `"important":true`) {
		t.Errorf("expected important:true, got %s", s)
	}
	if !strings.Contains(s, `"action":"[update:content:underlay:image|update:content:underlay:text]"`) {
		t.Errorf("expected compound action, got %s", s)
	}
}

// Item.Selection roundtrip: the pages0.json roundtrip already exercises
// this, but an explicit builder+marshal test here catches regressions in
// the field wiring.
func TestItem_WithSelection_Builder(t *testing.T) {
	it := content.NewItem().
		WithTitle("x").
		WithSelection(content.NewSelection().
			WithImportant(false).
			WithHeadline("focus").
			WithAction(action.Back))
	b, err := json.Marshal(it)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"selection":{`) {
		t.Errorf("expected selection object in output, got %s", s)
	}
	if !strings.Contains(s, `"important":false`) {
		t.Errorf("expected important:false preserved, got %s", s)
	}
}

// State objects (coming/running/over) may carry their own execute envelope,
// which must survive a round-trip.
func TestLive_StateExecuteRoundTrip(t *testing.T) {
	raw := `{"type":"schedule","coming":{"execute":{"action":"info:soon"}},"running":{"execute":{"action":"reload:content"}},"over":{"execute":{"action":"back"}}}`
	var l content.Live
	if err := json.Unmarshal([]byte(raw), &l); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := l.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	b, err := json.Marshal(&l)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got, want any
	_ = json.Unmarshal(b, &got)
	_ = json.Unmarshal([]byte(raw), &want)
	gb, _ := json.Marshal(got)
	wb, _ := json.Marshal(want)
	if string(gb) != string(wb) {
		t.Errorf("round-trip mismatch:\n got:  %s\n want: %s", gb, wb)
	}

	built := content.NewLive().WithType(content.LiveSchedule).
		WithOver(content.NewLiveState().WithExecute(action.Back, nil)).
		WithRunning(content.NewLiveRunningState().WithExecute(action.Reload, nil))
	b, _ = json.Marshal(built)
	for _, want := range []string{`"over":{"execute":{"action":"back"}}`, `"running":{"execute":{"action":`} {
		if !strings.Contains(string(b), want) {
			t.Errorf("expected %q in %s", want, b)
		}
	}
}

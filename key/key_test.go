package key_test

import (
	"encoding/json"
	"testing"

	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/key"
)

func TestKey_Constants(t *testing.T) {
	cases := map[string]key.Key{
		"play":               key.Play,
		"toggle_play_pause":  key.TogglePlayPause,
		"channel_up":         key.ChannelUp,
		"previous_track":     key.PreviousTrack,
		"bracket_open":       key.BracketOpen,
		"0":                  key.D0,
		"a":                  key.A,
	}
	for want, got := range cases {
		if string(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func TestCombine(t *testing.T) {
	if got := key.Combine(); got != "" {
		t.Errorf("empty combine should be empty, got %q", got)
	}
	if got := key.Combine(key.Play); got != key.Play {
		t.Errorf("single combine should pass through, got %q", got)
	}
	if got := key.Combine(key.Shift, key.CapsLock); string(got) != "shift|caps_lock" {
		t.Errorf("got %q", got)
	}
}

// Key values flow through the content.Item field with the correct JSON
// representation.
func TestItem_KeyFieldMarshal(t *testing.T) {
	it := content.NewItem().WithTitle("x").WithKey(key.Play)
	b, err := json.Marshal(it)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes_contains(b, `"key":"play"`) {
		t.Errorf(`expected "key":"play" in output, got %s`, b)
	}
}

// Casting a free-form string covers platform-specific bindings not in the
// constant set.
func TestKey_CustomString(t *testing.T) {
	it := content.NewItem().WithTitle("x").WithKey(key.Key("fast_forward"))
	b, _ := json.Marshal(it)
	if !bytes_contains(b, `"key":"fast_forward"`) {
		t.Errorf("custom key not round-tripped: %s", b)
	}
}

func bytes_contains(haystack []byte, needle string) bool {
	return indexString(haystack, needle) >= 0
}

func indexString(haystack []byte, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if string(haystack[i:i+len(needle)]) == needle {
			return i
		}
	}
	return -1
}

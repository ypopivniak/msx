package icon_test

import (
	"testing"

	"github.com/ypopivniak/msx/icon"
)

func TestTinted_BuildsCompositeString(t *testing.T) {
	got := icon.Tinted(icon.TintWhiteSoft, string(icon.Movie))
	if string(got) != "msx-white-soft:movie" {
		t.Errorf("got %q", got)
	}
}

func TestTinted_EmptyTintIsBareIcon(t *testing.T) {
	got := icon.Tinted("", string(icon.PlayArrow))
	if string(got) != "play-arrow" {
		t.Errorf("got %q", got)
	}
}

func TestTinted_AcceptsArbitraryTint(t *testing.T) {
	// MSX may recognize platform-specific or custom tints; the helper
	// should not gate on a fixed allowlist.
	got := icon.Tinted("custom-tint", "home")
	if string(got) != "custom-tint:home" {
		t.Errorf("got %q", got)
	}
}

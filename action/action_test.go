package action_test

import (
	"testing"

	"github.com/ypopivniak/msx/action"
)

func TestCompound_Shapes(t *testing.T) {
	cases := []struct {
		name string
		in   []action.Action
		want string
	}{
		{"empty", nil, "[]"},
		{"single-unbracketed", []action.Action{action.Back}, "back"},
		{
			name: "two-piped",
			in:   []action.Action{action.Invalidate("content"), action.ReloadScope("content")},
			want: "[invalidate:content|reload:content]",
		},
		{
			name: "three-piped",
			in:   []action.Action{action.Update("content", "x"), action.Focus("x"), action.Info("done")},
			want: "[update:content:x|focus:x|info:done]",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := string(action.Compound(tc.in...))
			if got != tc.want {
				t.Errorf("Compound(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestInvalidate_Variants(t *testing.T) {
	checks := map[string]action.Action{
		"invalidate:content":                 action.Invalidate("content"),
		"invalidate:menu:abc":                action.InvalidateItem("menu", "abc"),
		"invalidate:content:flag:home":       action.InvalidateFlag("content", "home"),
		"invalidate:panel:overlay:id1":       action.InvalidateOverlay("panel", "id1"),
		"invalidate:content:underlay:header": action.InvalidateUnderlay("content", "header"),
	}
	for want, got := range checks {
		if string(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func TestUpdate_Variants(t *testing.T) {
	if got := action.Update("content", "x"); string(got) != "update:content:x" {
		t.Errorf("Update got %q", got)
	}
	if got := action.UpdateOverlay("panel", "y"); string(got) != "update:panel:overlay:y" {
		t.Errorf("UpdateOverlay got %q", got)
	}
	if got := action.UpdateUnderlay("content", "z"); string(got) != "update:content:underlay:z" {
		t.Errorf("UpdateUnderlay got %q", got)
	}
}

func TestReloadScope(t *testing.T) {
	if got := action.ReloadScope(""); got != action.Reload {
		t.Errorf("empty scope should yield bare reload, got %q", got)
	}
	if got := action.ReloadScope("menu"); string(got) != "reload:menu" {
		t.Errorf("got %q", got)
	}
}

func TestFocus_Variants(t *testing.T) {
	if string(action.Focus("x")) != "focus:x" {
		t.Error("Focus wrong")
	}
	if string(action.FocusIndex(2)) != "focus:index:2" {
		t.Error("FocusIndex wrong")
	}
	if string(action.FocusExecute("y")) != "focus:execute:y" {
		t.Error("FocusExecute wrong")
	}
}

func TestPlayer_Constants(t *testing.T) {
	if string(action.PlayerPlayPause) != "player:play_pause" {
		t.Errorf("got %q", action.PlayerPlayPause)
	}
	if string(action.PlayerSeek("30")) != "player:seek:30" {
		t.Errorf("got %q", action.PlayerSeek("30"))
	}
	if string(action.PlayerGotoIndex(3)) != "player:goto:index:3" {
		t.Errorf("got %q", action.PlayerGotoIndex(3))
	}
}

func TestDelayLazy(t *testing.T) {
	if got := action.Delay("id1", 1500, action.Back); string(got) != "delay:id1:1500:back" {
		t.Errorf("got %q", got)
	}
	if got := action.Lazy(action.Reload); string(got) != "lazy:reload" {
		t.Errorf("got %q", got)
	}
}

// Reproduces the exact compound string seen in pages1.json.
func TestCompound_MatchesFixture(t *testing.T) {
	got := action.Compound(
		action.Invalidate("content"),
		action.ReloadScope("content"),
	)
	if string(got) != "[invalidate:content|reload:content]" {
		t.Errorf("got %q", got)
	}
}

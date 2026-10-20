package feed_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/feed"
)

func TestM3U_MinimalURL(t *testing.T) {
	c := feed.NewM3U("https://example.invalid/list.m3u")
	got, err := c.URLString()
	if err != nil {
		t.Fatalf("URLString: %v", err)
	}
	if !strings.HasPrefix(got, feed.DefaultM3UEndpoint+"?") {
		t.Errorf("expected default endpoint prefix, got %q", got)
	}

	// Parse the resulting query so the test doesn't depend on key order.
	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	q := u.Query()
	if q.Get("url") != "https://example.invalid/list.m3u" {
		t.Errorf("url param wrong: %q", got)
	}
	if len(q) != 1 {
		t.Errorf("expected only url param for minimal case, got %v", q)
	}
}

func TestM3U_FullParameters(t *testing.T) {
	c := feed.NewM3U("https://example.invalid/list.m3u").
		WithKind(feed.KindVideo).
		WithHeadline("My videos").
		WithBackground("https://example.invalid/bg.jpg").
		WithMode(feed.ModeExtended).
		WithContext(feed.ContextFavorites)

	got, err := c.URLString()
	if err != nil {
		t.Fatalf("URLString: %v", err)
	}
	u, _ := url.Parse(got)
	q := u.Query()
	checks := map[string]string{
		"url":        "https://example.invalid/list.m3u",
		"type":       "video",
		"headline":   "My videos",
		"background": "https://example.invalid/bg.jpg",
		"mode":       "extended",
		"context":    "favorites",
	}
	for k, want := range checks {
		if q.Get(k) != want {
			t.Errorf("%s = %q, want %q", k, q.Get(k), want)
		}
	}
}

func TestM3U_RejectsEmptyURL(t *testing.T) {
	_, err := feed.NewM3U("").URLString()
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "url") {
		t.Errorf("expected url error, got %v", err)
	}
}

func TestM3U_RejectsBadKind(t *testing.T) {
	_, err := feed.NewM3U("x").WithKind("banana").URLString()
	if err == nil || !strings.Contains(err.Error(), "type") {
		t.Errorf("expected type enum error, got %v", err)
	}
}

func TestPLS_UsesPLSEndpoint(t *testing.T) {
	c := feed.NewPLS("https://example.invalid/list.pls")
	got, _ := c.URLString()
	if !strings.HasPrefix(got, feed.DefaultPLSEndpoint+"?") {
		t.Errorf("expected pls endpoint, got %q", got)
	}
}

func TestConvert_CustomEndpoint(t *testing.T) {
	c := feed.NewM3U("x").WithEndpoint("https://mirror.example/m3u")
	got, _ := c.URLString()
	if !strings.HasPrefix(got, "https://mirror.example/m3u?") {
		t.Errorf("expected custom endpoint, got %q", got)
	}
}

func TestConvert_ContentAction(t *testing.T) {
	a, err := feed.NewM3U("https://example.invalid/list.m3u").
		WithKind(feed.KindAudio).
		ContentAction()
	if err != nil {
		t.Fatalf("ContentAction: %v", err)
	}
	if !strings.HasPrefix(string(a), "content:") {
		t.Errorf("expected content:<url> action, got %q", a)
	}
}

func TestMRSS_MinimalURL(t *testing.T) {
	got, err := feed.NewMRSS("https://example.invalid/feed.xml").URLString()
	if err != nil {
		t.Fatalf("URLString: %v", err)
	}
	want := "https://example.invalid/feed.xml@" + feed.DefaultMRSSPlugin
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMRSS_FullPipeParams(t *testing.T) {
	got, err := feed.NewMRSS("https://feed.invalid/x.xml").
		WithBanner("https://banner").
		WithBackground("https://bg").
		WithPluginURL("https://plugin.override").
		WithDictionary("https://dict").
		WithBrand("msx-red").
		URLString()
	if err != nil {
		t.Fatalf("URLString: %v", err)
	}
	want := "https://feed.invalid/x.xml@" + feed.DefaultMRSSPlugin +
		"?https://banner|https://bg|https://plugin.override|https://dict|msx-red"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

// Internal empty slots must be preserved so downstream positional
// parameters keep their meaning; trailing empties are trimmed.
func TestMRSS_InternalEmptyPreserved(t *testing.T) {
	got, _ := feed.NewMRSS("https://feed.invalid/x.xml").
		WithBackground("https://bg").
		URLString()
	want := "https://feed.invalid/x.xml@" + feed.DefaultMRSSPlugin + "?|https://bg"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestMRSS_TrailingEmptyTrimmed(t *testing.T) {
	got, _ := feed.NewMRSS("https://feed.invalid/x.xml").
		WithBanner("https://banner").
		URLString()
	want := "https://feed.invalid/x.xml@" + feed.DefaultMRSSPlugin + "?https://banner"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

func TestMRSS_RejectsEmptyFeed(t *testing.T) {
	_, err := feed.NewMRSS("").URLString()
	if err == nil || !strings.Contains(err.Error(), "feedURL") {
		t.Errorf("expected feedURL error, got %v", err)
	}
}

func TestMRSS_InteractionAction(t *testing.T) {
	a, err := feed.NewMRSS("https://feed.invalid/x.xml").InteractionAction()
	if err != nil {
		t.Fatalf("InteractionAction: %v", err)
	}
	want := "content:request:interaction:https://feed.invalid/x.xml@" + feed.DefaultMRSSPlugin
	if string(a) != want {
		t.Errorf("got %q\nwant %q", a, want)
	}
}

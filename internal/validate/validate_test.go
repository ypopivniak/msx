package validate

import "testing"

// Nested Prefix calls over joined errors must apply each segment exactly
// once to every field error.
func TestPrefix_NestedJoinAppliesEachSegmentOnce(t *testing.T) {
	item := Join(New("progress", "bad"), New("layout", "bad"))
	page := Join(Prefix("items[0]", item))
	root := Join(Prefix("pages[0]", page), New("type", "bad"))

	want := "pages[0].items[0].progress: bad\npages[0].items[0].layout: bad\ntype: bad"
	if got := root.Error(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestPrefix_EmptyPrefixLeavesPathUnchanged(t *testing.T) {
	err := Prefix("", Join(New("progress", "bad")))
	if got := err.Error(); got != "progress: bad" {
		t.Errorf("got %q, want %q", got, "progress: bad")
	}
}

func TestPrefix_EmptyPathTakesPrefix(t *testing.T) {
	err := Prefix("options", New("", "bad"))
	if got := err.Error(); got != "options: bad" {
		t.Errorf("got %q, want %q", got, "options: bad")
	}
}

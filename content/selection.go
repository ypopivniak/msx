package content

import (
	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/internal/validate"
)

// Selection models the MSX Selection Object used on a Content Item to
// customize focus behavior — typically triggering an update on overlay
// or underlay pages when an item receives focus.
//
// It lives in the content package (not a separate "selection" package)
// because it is attached to Content Item and embedding there would
// otherwise cause an import cycle.
type Selection struct {
	// Important pointer-valued so fixtures that explicitly write
	// "important": false or true round-trip faithfully.
	Important  *bool         `json:"important,omitempty"`
	Headline   string        `json:"headline,omitempty"`
	Background string        `json:"background,omitempty"`
	Action     action.Action `json:"action,omitempty"`
	Data       any           `json:"data,omitempty"`
}

// NewSelection returns a fresh Selection ready for chained configuration.
func NewSelection() *Selection { return &Selection{} }

func (s *Selection) WithImportant(v bool) *Selection  { s.Important = &v; return s }
func (s *Selection) WithHeadline(v string) *Selection { s.Headline = v; return s }
func (s *Selection) WithBackground(url string) *Selection {
	s.Background = url
	return s
}
func (s *Selection) WithAction(a action.Action) *Selection { s.Action = a; return s }
func (s *Selection) WithData(d any) *Selection            { s.Data = d; return s }

// Validate is a no-op at present — the Selection Object has no strict
// per-field constraints in the MSX spec beyond JSON type safety, which
// the Go struct already enforces. The method is kept so callers can
// treat Selection uniformly with the rest of the library.
func (s *Selection) Validate() error {
	if s == nil {
		return nil
	}
	// Placeholder for future cross-field rules; returning validate.Join
	// with no args keeps the import used and documents the intent.
	return validate.Join()
}

package menu

import (
	"encoding/json"

	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/internal/validate"
)

// ItemData holds the polymorphic `data` property of a Menu Item, which MSX
// accepts as either a URL/request-action string (the item loads content from
// that address) or an inline Content Root Object.
//
// Exactly one of Ref or Inline should be set. The MarshalJSON picks between
// them; Validate enforces the exclusivity so malformed trees are caught at
// Build() time rather than producing ambiguous JSON.
type ItemData struct {
	// Ref is a URL or MSX request action string like "content:https://…".
	Ref string
	// Inline is an embedded Content Root, used when the item's content is
	// described directly instead of fetched.
	Inline *content.Content
}

// NewItemDataRef constructs an ItemData referencing an external URL/action.
func NewItemDataRef(ref string) *ItemData { return &ItemData{Ref: ref} }

// NewItemDataInline constructs an ItemData with an inline Content Root.
func NewItemDataInline(c *content.Content) *ItemData { return &ItemData{Inline: c} }

// MarshalJSON writes Ref as a bare string or Inline as an object. An empty
// ItemData marshals to null so the containing omitempty tag drops it; that
// case should not arise because ItemData is always held as a pointer.
func (d *ItemData) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	if d.Ref != "" {
		return json.Marshal(d.Ref)
	}
	if d.Inline != nil {
		return json.Marshal(d.Inline)
	}
	return []byte("null"), nil
}

// UnmarshalJSON accepts either a JSON string (stored as Ref) or a JSON
// object (decoded into Inline), mirroring the polymorphic MSX spec.
func (d *ItemData) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		d.Ref = s
		return nil
	}
	var c content.Content
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	d.Inline = &c
	return nil
}

// Validate ensures at most one representation is set and recurses into
// the inline content if present.
func (d *ItemData) Validate() error {
	if d == nil {
		return nil
	}
	if d.Ref != "" && d.Inline != nil {
		return validate.New("", "ref and inline content are mutually exclusive")
	}
	if d.Ref == "" && d.Inline == nil {
		return validate.New("", "either ref or inline content must be set")
	}
	if d.Inline != nil {
		return d.Inline.Validate()
	}
	return nil
}

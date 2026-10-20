// Package content models the MSX Content Root / Page / Item tree using a
// fluent builder. Calling Build() on the root performs strict validation
// and returns a pointer suitable for json.Marshal.
package content

import (
	"encoding/json"

	"github.com/ypopivniak/msx/internal/validate"
)

// Type is the top-level content layout direction.
type Type string

const (
	// TypePages is the default: items are distributed across discrete pages
	// that occupy the full screen.
	TypePages Type = "pages"
	// TypeList lays items out as a single scrollable list.
	TypeList Type = "list"
)

var allTypes = []Type{TypePages, TypeList}

// Preload controls whether adjacent pages are fetched in advance.
type Preload string

const (
	PreloadNone Preload = "none"
	PreloadNext Preload = "next"
	PreloadPrev Preload = "prev"
	PreloadFull Preload = "full"
)

var allPreloads = []Preload{PreloadNone, PreloadNext, PreloadPrev, PreloadFull}

// ItemType is the visual style of a content item.
type ItemType string

const (
	ItemTypeDefault  ItemType = "default"
	ItemTypeTeaser   ItemType = "teaser"
	ItemTypeButton   ItemType = "button"
	ItemTypeSeparate ItemType = "separate"
	ItemTypeSpace    ItemType = "space"
	ItemTypeControl  ItemType = "control"
)

var allItemTypes = []ItemType{
	ItemTypeDefault, ItemTypeTeaser, ItemTypeButton,
	ItemTypeSeparate, ItemTypeSpace, ItemTypeControl,
}

// ImageFiller controls how the background image scales inside the item.
type ImageFiller string

const (
	FillerDefault      ImageFiller = "default"
	FillerWidth        ImageFiller = "width"
	FillerWidthTop     ImageFiller = "width-top"
	FillerWidthCenter  ImageFiller = "width-center"
	FillerWidthBottom  ImageFiller = "width-bottom"
	FillerHeight       ImageFiller = "height"
	FillerHeightLeft   ImageFiller = "height-left"
	FillerHeightCenter ImageFiller = "height-center"
	FillerHeightRight  ImageFiller = "height-right"
	FillerFit          ImageFiller = "fit"
	FillerCover        ImageFiller = "cover"
	FillerSmart        ImageFiller = "smart"
)

var allFillers = []ImageFiller{
	FillerDefault, FillerWidth, FillerWidthTop, FillerWidthCenter, FillerWidthBottom,
	FillerHeight, FillerHeightLeft, FillerHeightCenter, FillerHeightRight,
	FillerFit, FillerCover, FillerSmart,
}

// Break forces a line/page boundary after an item.
type Break string

const (
	BreakNone Break = ""
	BreakLine Break = "line"
	BreakPage Break = "page"
	// BreakContext ends the current context; attach an id with
	// BreakContextID(id) if needed.
	BreakContext Break = "context"
)

// BreakContextID builds "context:<id>" for the `break` field.
func BreakContextID(id string) Break { return Break("context:" + id) }

// Tri is a tri-state value used by properties documented as "number|boolean"
// with the specific values {0 or false, 1 or true, 2}. We always emit a
// number so output is stable regardless of how the value was set.
type Tri int

const (
	TriOff  Tri = 0
	TriOn   Tri = 1
	TriAuto Tri = 2
)

// UnmarshalJSON accepts the boolean form (false → 0, true → 1) as well as a
// plain number. Range checking is left to the owning object's Validate.
func (t *Tri) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*t = TriOn
		} else {
			*t = TriOff
		}
		return nil
	}
	var n int
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*t = Tri(n)
	return nil
}

// validateOneOf is the shared "value must be in allowed set" helper used by
// all enum validations in this package.
func validateOneOf[T comparable](path string, value T, allowed []T) error {
	if len(allowed) == 0 {
		return nil
	}
	return validate.OneOf(path, value, allowed...)
}

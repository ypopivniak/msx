// Package menu models the MSX Menu Root Object and Menu Item Object using
// the same fluent-builder, strict-validation pattern as the content package.
//
// A Menu Root is the top-level navigation document: a vertical list of items
// on the left side of the screen, each of which can open an inline or
// referenced Content Root. The field shapes mirror the MSX wiki one-for-one.
package menu

import "github.com/ypopivniak/msx/content"

// Style controls the overall menu display treatment.
type Style string

const (
	StyleDefault          Style = "default"
	StyleFlat             Style = "flat"
	StyleFlatSeparator    Style = "flat-separator"
	StyleOverlay          Style = "overlay"
	StyleOverlaySeparator Style = "overlay-separator"
)

var allStyles = []Style{
	StyleDefault, StyleFlat, StyleFlatSeparator,
	StyleOverlay, StyleOverlaySeparator,
}

// LogoSize is the fixed pair of allowed values for the `logoSize` property.
type LogoSize string

const (
	LogoSmall LogoSize = "small"
	LogoLarge LogoSize = "large"
)

var allLogoSizes = []LogoSize{LogoSmall, LogoLarge}

// ItemType is the visual style of a menu item.
type ItemType string

const (
	ItemTypeDefault   ItemType = "default"
	ItemTypeSeparator ItemType = "separator"
	ItemTypeSettings  ItemType = "settings"
)

var allItemTypes = []ItemType{ItemTypeDefault, ItemTypeSeparator, ItemTypeSettings}

// Tri aliases content.Tri so callers have one canonical tri-state numeric
// type across MSX objects without an extra import.
type Tri = content.Tri

// Re-export the Tri constants for convenience.
const (
	TriOff  = content.TriOff
	TriOn   = content.TriOn
	TriAuto = content.TriAuto
)

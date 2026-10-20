// Package color defines MSX color names.
//
// MSX identifies colors by name (e.g. "msx-glass", "msx-yellow"). Because the
// palette is open-ended (plugins may define custom colors, and base names
// accept "-transparent"/"-soft"/tinted suffixes) this package exposes a
// typed string alias with the most common constants rather than a closed enum.
// Any non-empty string is accepted; callers wanting stricter guarantees can
// compare against the provided constants.
package color

// Color is an MSX color identifier. Empty means "unset" and is omitted from
// the resulting JSON by the containing struct's omitempty tag.
type Color string

func (c Color) String() string { return string(c) }

// Core palette. See https://msx.benzac.de/wiki/index.php?title=Colors
const (
	None Color = ""

	MSXWhite  Color = "msx-white"
	MSXBlack  Color = "msx-black"
	MSXGlass  Color = "msx-glass"
	MSXGray   Color = "msx-gray"
	MSXRed    Color = "msx-red"
	MSXGreen  Color = "msx-green"
	MSXBlue   Color = "msx-blue"
	MSXYellow Color = "msx-yellow"
	MSXOrange Color = "msx-orange"
	MSXViolet Color = "msx-violet"
	MSXTeal   Color = "msx-teal"
	MSXCyan   Color = "msx-cyan"
	MSXPink   Color = "msx-pink"
	MSXBrown  Color = "msx-brown"
)

// Package icon defines MSX icon identifiers.
//
// MSX uses Material icon names (e.g. "home", "play-arrow") plus a small set of
// MSX-specific icons ("msx-white:home", etc.). The catalog is large and
// open-ended, so this package exposes a typed string alias with a handful of
// commonly used constants rather than enumerating the full set.
package icon

// Icon is an MSX icon identifier. The empty value means "unset".
type Icon string

func (i Icon) String() string { return string(i) }

// Tinted returns "<tint>:<name>" — the MSX composite form used to
// prefix a base icon with a color hint such as "msx-white-soft:movie".
// The tint is typically a Color constant from the color package, but any
// identifier MSX recognizes (including "msx-white-soft", "msx-white",
// etc.) is accepted.
func Tinted(tint, name string) Icon {
	if tint == "" {
		return Icon(name)
	}
	return Icon(tint + ":" + name)
}

// Common icon-tint modifiers used across MSX reference content.
const (
	TintWhite     = "msx-white"
	TintWhiteSoft = "msx-white-soft"
	TintBlack     = "msx-black"
	TintBlackSoft = "msx-black-soft"
)

// Frequently used icons. See https://msx.benzac.de/wiki/index.php?title=Icons
const (
	None Icon = ""

	Home       Icon = "home"
	Search     Icon = "search"
	Settings   Icon = "settings"
	Info       Icon = "info"
	Warning    Icon = "warning"
	Error      Icon = "error"
	Check      Icon = "check"
	Close      Icon = "close"
	Refresh    Icon = "refresh"
	Menu       Icon = "menu"
	MoreVert   Icon = "more-vert"
	MoreHoriz  Icon = "more-horiz"
	PlayArrow  Icon = "play-arrow"
	Pause      Icon = "pause"
	Stop       Icon = "stop"
	SkipNext   Icon = "skip-next"
	SkipPrev   Icon = "skip-previous"
	VolumeUp   Icon = "volume-up"
	VolumeOff  Icon = "volume-off"
	Folder     Icon = "folder"
	FolderOpen Icon = "folder-open"
	Movie      Icon = "movie"
	Tv         Icon = "tv"
	MusicNote  Icon = "music-note"
	Image      Icon = "image"
	Star       Icon = "star"
	Favorite   Icon = "favorite"
)

// IconSize is a valid value for the `iconSize` property of a content item.
type Size string

const (
	SizeSmall      Size = "small"
	SizeMedium     Size = "medium"
	SizeLarge      Size = "large"
	SizeExtraLarge Size = "extra-large"
)

// AllSizes lists every valid IconSize value, useful for validation.
var AllSizes = []Size{SizeSmall, SizeMedium, SizeLarge, SizeExtraLarge}

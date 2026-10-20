// Package key defines typed identifiers for the MSX `key` property on
// content items, which binds a remote-control button or keyboard key
// (or a pipe-separated group of them) to an item's action.
//
// See https://msx.benzac.de/wiki/index.php?title=Key_Property.
package key

import "strings"

// Key is a typed string alias. Any string is accepted (MSX may recognize
// unlisted bindings on specific platforms), so call sites can fall back
// to Key("custom") when the constant set doesn't cover a variant.
type Key string

// String returns the raw MSX key identifier.
func (k Key) String() string { return string(k) }

// Navigation keys.
const (
	Home   Key = "home"
	End    Key = "end"
	Insert Key = "insert"
	Delete Key = "delete"
)

// Playback controls.
const (
	Play            Key = "play"
	Pause           Key = "pause"
	Stop            Key = "stop"
	TogglePlayPause Key = "toggle_play_pause"
	Forward         Key = "forward"
	Rewind          Key = "rewind"
	NextTrack       Key = "next_track"
	PreviousTrack   Key = "previous_track"
	Record          Key = "record"
	Restart         Key = "restart"
)

// Channel & menu controls.
const (
	ChannelUp   Key = "channel_up"
	ChannelDown Key = "channel_down"
	Guide       Key = "guide"
	Info        Key = "info"
	Clear       Key = "clear"
	Settings    Key = "settings"
	Search      Key = "search"
	ChannelList Key = "channel_list"
)

// Color / function buttons.
const (
	Red    Key = "red"
	Green  Key = "green"
	Yellow Key = "yellow"
	Blue   Key = "blue"
)

// Digits.
const (
	D0 Key = "0"
	D1 Key = "1"
	D2 Key = "2"
	D3 Key = "3"
	D4 Key = "4"
	D5 Key = "5"
	D6 Key = "6"
	D7 Key = "7"
	D8 Key = "8"
	D9 Key = "9"
)

// Letters (lowercase only — MSX uses unshifted forms for bindings).
const (
	A Key = "a"
	B Key = "b"
	C Key = "c"
	D Key = "d"
	E Key = "e"
	F Key = "f"
	G Key = "g"
	H Key = "h"
	I Key = "i"
	J Key = "j"
	K Key = "k"
	L Key = "l"
	M Key = "m"
	N Key = "n"
	O Key = "o"
	P Key = "p"
	Q Key = "q"
	R Key = "r"
	S Key = "s"
	T Key = "t"
	U Key = "u"
	V Key = "v"
	W Key = "w"
	X Key = "x"
	Y Key = "y"
	Z Key = "z"
)

// Symbols & punctuation.
const (
	Space        Key = "space"
	Tab          Key = "tab"
	Period       Key = "period"
	Comma        Key = "comma"
	Semicolon    Key = "semicolon"
	Dash         Key = "dash"
	Equal        Key = "equal"
	Quote        Key = "quote"
	Accent       Key = "accent"
	Slash        Key = "slash"
	Backslash    Key = "backslash"
	BracketOpen  Key = "bracket_open"
	BracketClose Key = "bracket_close"
)

// Modifiers.
const (
	Shift    Key = "shift"
	CapsLock Key = "caps_lock"
)

// Combine joins multiple keys into the MSX pipe-separated form so any of
// them triggers the bound action. Calling Combine with a single key
// returns it unchanged; with none returns the empty Key.
func Combine(keys ...Key) Key {
	switch len(keys) {
	case 0:
		return ""
	case 1:
		return keys[0]
	}
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = string(k)
	}
	return Key(strings.Join(parts, "|"))
}

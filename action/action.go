// Package action builds MSX action strings.
//
// MSX actions are strings of the form "<prefix>:<value>" (e.g.
// "video:https://example.com/stream.m3u8") or a bare keyword ("back",
// "execute", "reload"). This package provides typed constructors so that
// callers do not build action strings by concatenation.
//
// The catalog is large — the MSX wiki lists hundreds of internal action
// variants. We expose the commonly-used families as dedicated helpers and
// leave Prefixed() as the escape hatch for anything not modelled.
package action

import (
	"strings"
)

// Action is a raw MSX action string. Use the constructors below to produce
// well-formed values, or cast directly for advanced actions not modelled here.
type Action string

func (a Action) String() string { return string(a) }

// Plain returns the action as-is (e.g. "back", "execute"). Useful for the
// bare keyword actions that do not take a payload.
func Plain(name string) Action { return Action(name) }

// Prefixed returns "<prefix>:<value>", panicking on empty prefix. This is the
// lowest-level escape hatch for action kinds not covered by the helpers below.
func Prefixed(prefix, value string) Action {
	if prefix == "" {
		panic("action: empty prefix")
	}
	return Action(prefix + ":" + value)
}

// Common prefix constructors. The set is intentionally small; add more as the
// library grows. See https://msx.benzac.de/wiki/index.php?title=Actions.

func Video(url string) Action   { return Prefixed("video", url) }
func Audio(url string) Action   { return Prefixed("audio", url) }
func Image(url string) Action   { return Prefixed("image", url) }
func Content(url string) Action { return Prefixed("content", url) }
func Panel(url string) Action   { return Prefixed("panel", url) }
func Link(url string) Action    { return Prefixed("link", url) }
func Info(text string) Action   { return Prefixed("info", text) }
func HTTP(url string) Action    { return Prefixed("http", url) }
func HTTPS(url string) Action   { return Prefixed("https", url) }

// Bare-keyword actions.
const (
	Back       Action = "back"
	Reload     Action = "reload"
	Execute    Action = "execute"
	Cancel     Action = "cancel"
	Home       Action = "home"
	Main       Action = "main"
	Close      Action = "close"
	Cleanup    Action = "cleanup"
	Exit       Action = "exit"
	Restart    Action = "restart"
	Fullscreen Action = "fullscreen"
	Awake      Action = "awake"
	Sleep      Action = "sleep"
	Leave      Action = "leave"
	Eject      Action = "eject"
	Start      Action = "start"
	Selection  Action = "selection"
	Live       Action = "live"
	Settings   Action = "settings"
	Options    Action = "options"
	Playback   Action = "playback"
	Player     Action = "player"
	Data       Action = "data"
	Noop       Action = "[]"
	None       Action = ""
)

// Prefix returns the part before the first ":" or the full string if none is
// present. Useful when validating an action against an expected family.
func (a Action) Prefix() string {
	s := string(a)
	if i := strings.IndexByte(s, ':'); i >= 0 {
		return s[:i]
	}
	return s
}

// Compound composes multiple actions into the MSX pipe-separated form
// "[a|b|c]", which MSX executes sequentially. A single action is returned
// as-is (no brackets) so it round-trips identically to its bare form. An
// empty call returns Noop ("[]").
func Compound(actions ...Action) Action {
	switch len(actions) {
	case 0:
		return Noop
	case 1:
		return actions[0]
	}
	parts := make([]string, len(actions))
	for i, a := range actions {
		parts[i] = string(a)
	}
	return Action("[" + strings.Join(parts, "|") + "]")
}

// ---- Invalidate family ----

// Invalidate produces "invalidate:<scope>" where scope is one of
// "content", "menu", or "panel".
func Invalidate(scope string) Action { return Prefixed("invalidate", scope) }

// InvalidateItem produces "invalidate:<scope>:<id>".
func InvalidateItem(scope, id string) Action {
	return Prefixed("invalidate", scope+":"+id)
}

// InvalidateFlag produces "invalidate:<scope>:flag:<flag>".
func InvalidateFlag(scope, flag string) Action {
	return Prefixed("invalidate", scope+":flag:"+flag)
}

// InvalidateOverlay produces "invalidate:<scope>:overlay:<id>".
func InvalidateOverlay(scope, id string) Action {
	return Prefixed("invalidate", scope+":overlay:"+id)
}

// InvalidateUnderlay produces "invalidate:<scope>:underlay:<id>".
func InvalidateUnderlay(scope, id string) Action {
	return Prefixed("invalidate", scope+":underlay:"+id)
}

// ---- Update family ----

// Update produces "update:<scope>:<id>" for the common case of updating
// a content/menu/panel item by id.
func Update(scope, id string) Action { return Prefixed("update", scope+":"+id) }

// UpdateOverlay / UpdateUnderlay target the overlay and underlay pages.
func UpdateOverlay(scope, id string) Action  { return Prefixed("update", scope+":overlay:"+id) }
func UpdateUnderlay(scope, id string) Action { return Prefixed("update", scope+":underlay:"+id) }

// ---- Reload family ----

// ReloadScope produces "reload:<scope>". Passing "" returns bare Reload.
func ReloadScope(scope string) Action {
	if scope == "" {
		return Reload
	}
	return Prefixed("reload", scope)
}

// ---- Focus family ----

// Focus produces "focus:<id>".
func Focus(id string) Action { return Prefixed("focus", id) }

// FocusIndex produces "focus:index:<i>".
func FocusIndex(i int) Action { return Prefixed("focus", "index:"+itoa(i)) }

// FocusExecute produces "focus:execute:<id>".
func FocusExecute(id string) Action { return Prefixed("focus", "execute:"+id) }

// FocusControl produces "focus:control:<name>".
func FocusControl(name string) Action { return Prefixed("focus", "control:"+name) }

// FocusClear produces "focus:clear:<flag>".
func FocusClear(flag string) Action { return Prefixed("focus", "clear:"+flag) }

// ---- Message family ----

func Success(text string) Action { return Prefixed("success", text) }
func Warn(text string) Action    { return Prefixed("warn", text) }
func Error(text string) Action   { return Prefixed("error", text) }

// ---- Logger family ----

func LoggerInfo(text string) Action    { return Prefixed("logger", "info:"+text) }
func LoggerError(text string) Action   { return Prefixed("logger", "error:"+text) }
func LoggerWarn(text string) Action    { return Prefixed("logger", "warn:"+text) }
func LoggerDebug(text string) Action   { return Prefixed("logger", "debug:"+text) }
func LoggerSuccess(text string) Action { return Prefixed("logger", "success:"+text) }

// ---- Link family ----

// LinkWindow opens the URL in a new browser window.
func LinkWindow(url string) Action { return Prefixed("link", "window:"+url) }

// LinkReplace opens the URL replacing the current window.
func LinkReplace(url string) Action { return Prefixed("link", "replace:"+url) }

// LinkValidate opens the URL after validation.
func LinkValidate(url string) Action { return Prefixed("link", "validate:"+url) }

// ---- Dialog / notification ----

func Dialog(id string) Action        { return Prefixed("dialog", id) }
func Notification(msg string) Action { return Prefixed("notification", msg) }

// ---- Player family (subset of the most common controls) ----

const (
	PlayerPlay      Action = "player:play"
	PlayerPause     Action = "player:pause"
	PlayerPlayPause Action = "player:play_pause"
	PlayerStop      Action = "player:stop"
	PlayerEject     Action = "player:eject"
	PlayerNext      Action = "player:next"
	PlayerPrev      Action = "player:prev"
	PlayerForward   Action = "player:forward"
	PlayerRewind    Action = "player:rewind"
	PlayerRestart   Action = "player:restart"
	PlayerShow      Action = "player:show"
	PlayerHide      Action = "player:hide"
)

// PlayerSeek produces "player:seek:<value>" where value is seconds or a
// percentage depending on MSX player semantics.
func PlayerSeek(value string) Action { return Prefixed("player", "seek:"+value) }

// PlayerGoto produces "player:goto:<target>" for "next"/"prev"/"first"/"last"
// or an item id.
func PlayerGoto(target string) Action { return Prefixed("player", "goto:"+target) }

// PlayerGotoIndex produces "player:goto:index:<i>".
func PlayerGotoIndex(i int) Action { return Prefixed("player", "goto:index:"+itoa(i)) }

// ---- Delay / lazy / event ----

// Delay produces "delay:<id>:<ms>:<action>".
func Delay(id string, ms int, a Action) Action {
	return Prefixed("delay", id+":"+itoa(ms)+":"+string(a))
}

// DelayCancel produces "delay:<id>:cancel".
func DelayCancel(id string) Action { return Prefixed("delay", id+":cancel") }

// Lazy wraps an action so it runs once pending animations complete.
func Lazy(a Action) Action { return Prefixed("lazy", string(a)) }

// Event triggers a custom event by id.
func Event(id string) Action { return Prefixed("event", id) }

// ---- Settings family (common setters) ----

func SettingsAbout() Action               { return Action("settings:about") }
func SettingsAnimate(mode int) Action     { return Prefixed("settings", "animate:"+itoa(mode)) }
func SettingsZoom(value string) Action    { return Prefixed("settings", "zoom:"+value) }

// ---- Menu / Content / Panel load variants ----

// Menu produces "menu:<url>".
func Menu(url string) Action { return Prefixed("menu", url) }

// MenuData produces "menu:data".
func MenuData() Action { return Action("menu:data") }

// ContentData produces "content:data".
func ContentData() Action { return Action("content:data") }

// PanelData produces "panel:data".
func PanelData() Action { return Action("panel:data") }

// Playlist / Slideshow producers.
func Playlist(url string) Action      { return Prefixed("playlist", url) }
func PlaylistData() Action            { return Action("playlist:data") }
func Slideshow(url string) Action     { return Prefixed("slideshow", url) }
func SlideshowData() Action           { return Action("slideshow:data") }

// itoa is a tiny helper so this file doesn't import strconv just to format
// integers in action strings.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

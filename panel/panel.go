// Package panel provides builder-style helpers for the MSX options panel
// (the `options` property found on menu, content root, page, and item).
//
// An MSX panel is not a distinct object type in the wiki: it is simply a
// Content Page Object or Content Root Object displayed in an overlay. The
// polymorphic wrapper lives in the content package as content.Panel; this
// package exposes thin constructors and convenience helpers so panel code
// reads as panel.FromPage(...), panel.New()... etc. rather than forcing
// callers to import content solely for the wrapper type.
package panel

import (
	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/icon"
)

// Panel re-exports the polymorphic wrapper so client code can use
// panel.Panel without importing content.
type Panel = content.Panel

// FromPage wraps a Content Page as an options panel.
func FromPage(p *content.Page) *Panel { return content.PanelFromPage(p) }

// FromContent wraps a Content Root as an options panel.
func FromContent(c *content.Content) *Panel { return content.PanelFromContent(c) }

// NewPage returns a Content Page pre-configured for typical options-panel
// use: transparent background and compressed layout. Callers override any
// field as usual via the Page's With* methods.
func NewPage() *content.Page {
	return content.NewPage().
		WithTransparent(content.TriOn).
		WithCompress(true)
}

// Button is a convenience for defining a single-action button item inside
// a panel page: a labelled default-type item that executes the given
// action when focused.
func Button(label string, a action.Action) *content.Item {
	return content.NewItem().
		WithType(content.ItemTypeButton).
		WithLabel(label).
		WithAction(a)
}

// MessageOption controls the tone of Message / Info / Warn / Success / Error
// panels built by this package. The value selects both the icon and the
// wrapper color so call sites don't hardcode visual decisions.
type MessageOption struct {
	Icon  icon.Icon
	Color color.Color
}

// Preset message tones.
var (
	ToneInfo    = MessageOption{Icon: icon.Info, Color: color.MSXBlue}
	ToneWarn    = MessageOption{Icon: icon.Warning, Color: color.MSXYellow}
	ToneError   = MessageOption{Icon: icon.Error, Color: color.MSXRed}
	ToneSuccess = MessageOption{Icon: icon.Check, Color: color.MSXGreen}
)

// Message builds a panel that displays a headline + body message with a
// single dismiss button. Use Info / Warn / Success / Error for the common
// tones; use Message directly to customize the icon and color.
func Message(headline, body string, tone MessageOption) *Panel {
	page := NewPage().
		AddItem(content.NewItem().
			WithType(content.ItemTypeDefault).
			WithLayoutXYWH(0, 0, 8, 4).
			WithIcon(iconOf(tone)).
			WithColor(colorOf(tone)).
			WithHeadline(headline).
			WithText(body)).
		AddItem(Button("OK", action.Back).WithLayoutXYWH(0, 4, 8, 1))
	return FromPage(page)
}

// Info / Warn / Success / Error are convenience wrappers around Message
// with the matching tone preset.
func Info(headline, body string) *Panel    { return Message(headline, body, ToneInfo) }
func Warn(headline, body string) *Panel    { return Message(headline, body, ToneWarn) }
func Success(headline, body string) *Panel { return Message(headline, body, ToneSuccess) }
func Error(headline, body string) *Panel   { return Message(headline, body, ToneError) }

// Confirm builds a two-button panel for yes/no prompts. `onConfirm` is the
// action run when the user picks the accept button; `onCancel` is run (or
// just back-navigates) for the other choice — pass action.Back for the
// default behavior.
func Confirm(headline, body string, onConfirm, onCancel action.Action) *Panel {
	accept := onConfirm
	cancel := onCancel
	if cancel == "" {
		cancel = action.Back
	}
	page := NewPage().
		AddItem(content.NewItem().
			WithType(content.ItemTypeDefault).
			WithLayoutXYWH(0, 0, 8, 3).
			WithHeadline(headline).
			WithText(body)).
		AddItem(Button("Yes", accept).WithLayoutXYWH(0, 3, 4, 2).WithFocus(true)).
		AddItem(Button("No", cancel).WithLayoutXYWH(4, 3, 4, 2))
	return FromPage(page)
}

func iconOf(t MessageOption) icon.Icon {
	if t.Icon != "" {
		return t.Icon
	}
	return ToneInfo.Icon
}

func colorOf(t MessageOption) color.Color {
	if t.Color != "" {
		return t.Color
	}
	return ToneInfo.Color
}

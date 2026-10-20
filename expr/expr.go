// Package expr builds MSX inline-expression strings safely.
//
// MSX string fields (headline, text, title, …) accept a small template
// language: "{txt:COLOR:TEXT}", "{ico:COLOR:ICON}", "{br}" etc. These helpers
// let callers compose those expressions without concatenating braces and
// colons by hand, and without risking malformed output that MSX will emit
// verbatim as plain text. See:
//
//	https://msx.benzac.de/wiki/index.php?title=Inline_Expressions
package expr

import (
	"fmt"
	"strings"

	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/icon"
)

// Zero-argument tokens. These are plain strings so they compose trivially
// via + / fmt.Sprintf without wrapping types.
const (
	Space     = "{sp}"
	Tab       = "{tb}"
	LineBreak = "{br}"
)

// Text emits "{txt:COLOR:TEXT}" when a color is provided, "{txt:TEXT}" (a
// form accepted by MSX for unstyled wrapping) when color is empty.
func Text(c color.Color, text string) string {
	if c == "" {
		return "{txt:" + text + "}"
	}
	return "{txt:" + string(c) + ":" + text + "}"
}

// Icon emits "{ico:COLOR:ICON}" or "{ico:ICON}" when color is empty.
func Icon(c color.Color, ic icon.Icon) string {
	if c == "" {
		return "{ico:" + string(ic) + "}"
	}
	return "{ico:" + string(c) + ":" + string(ic) + "}"
}

// Char emits "{chr:CHAR}" where CHAR is an HTML entity name (e.g. "copy")
// or a numeric entity like "#169".
func Char(name string) string { return "{chr:" + name + "}" }

// CharColored emits "{chr:COLOR:CHAR}".
func CharColored(c color.Color, name string) string {
	return "{chr:" + string(c) + ":" + name + "}"
}

// Col emits "{col:COLOR}" — sets the color for subsequent plain text in
// the same string field. See MSX's "Color Inline Expressions".
func Col(c color.Color) string { return "{col:" + string(c) + "}" }

// Context emits "{context:PROPERTY}" — resolves to the given context
// property value at render time (used inside templates).
func Context(property string) string { return "{context:" + property + "}" }

// Dict emits "{dict:KEY}" — looks up a dictionary translation.
func Dict(key string) string { return "{dict:" + key + "}" }

// Number emits "{num:VALUE:format}" which applies MSX's default number
// formatting (thousands separator applied above 999).
func Number(value any) string { return "{num:" + fmt.Sprint(value) + ":format}" }

// NumberFormat emits "{num:VALUE:format:FORMAT}" where FORMAT is the MSX
// number-format string (e.g. "2" for two decimals, " ,24" for custom
// separators and digits, per the wiki).
func NumberFormat(value any, format string) string {
	return "{num:" + fmt.Sprint(value) + ":format:" + format + "}"
}

// Date emits "{num:VALUE:date:FORMAT}". VALUE may be a millisecond
// timestamp or the literal keyword "now" for the current time.
func Date(value any, format string) string {
	return "{num:" + fmt.Sprint(value) + ":date:" + format + "}"
}

// Duration emits "{num:VALUE:duration:FORMAT}" — formats a millisecond
// duration as text according to FORMAT (e.g. "text:dhms").
func Duration(value any, format string) string {
	return "{num:" + fmt.Sprint(value) + ":duration:" + format + "}"
}

// Join concatenates expression parts with no separator. Use it when you
// want to compose a longer string without building a Builder.
func Join(parts ...string) string { return strings.Join(parts, "") }

// Builder is a fluent concatenator for inline expressions. It exists so
// that call sites read top-to-bottom instead of nesting fmt.Sprintf calls.
type Builder struct {
	sb strings.Builder
}

// New returns an empty Builder.
func New() *Builder { return &Builder{} }

// Raw appends a literal string (no escaping — MSX accepts plain text
// interleaved with inline-expression tokens).
func (b *Builder) Raw(s string) *Builder { b.sb.WriteString(s); return b }

// Text / Icon / Char etc. mirror the package-level constructors.
func (b *Builder) Text(c color.Color, text string) *Builder {
	b.sb.WriteString(Text(c, text))
	return b
}
func (b *Builder) Icon(c color.Color, ic icon.Icon) *Builder {
	b.sb.WriteString(Icon(c, ic))
	return b
}
func (b *Builder) Char(name string) *Builder { b.sb.WriteString(Char(name)); return b }
func (b *Builder) CharColored(c color.Color, name string) *Builder {
	b.sb.WriteString(CharColored(c, name))
	return b
}
func (b *Builder) Col(c color.Color) *Builder   { b.sb.WriteString(Col(c)); return b }
func (b *Builder) Context(prop string) *Builder { b.sb.WriteString(Context(prop)); return b }
func (b *Builder) Dict(key string) *Builder     { b.sb.WriteString(Dict(key)); return b }
func (b *Builder) Number(v any) *Builder        { b.sb.WriteString(Number(v)); return b }
func (b *Builder) NumberFormat(v any, f string) *Builder {
	b.sb.WriteString(NumberFormat(v, f))
	return b
}
func (b *Builder) Date(v any, format string) *Builder {
	b.sb.WriteString(Date(v, format))
	return b
}
func (b *Builder) Duration(v any, format string) *Builder {
	b.sb.WriteString(Duration(v, format))
	return b
}

// Br / Sp / Tb append zero-argument tokens.
func (b *Builder) Br() *Builder { b.sb.WriteString(LineBreak); return b }
func (b *Builder) Sp() *Builder { b.sb.WriteString(Space); return b }
func (b *Builder) Tb() *Builder { b.sb.WriteString(Tab); return b }

// String returns the assembled expression.
func (b *Builder) String() string { return b.sb.String() }

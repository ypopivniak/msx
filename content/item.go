package content

import (
	"encoding/json"
	"regexp"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/geom"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/internal/validate"
	"github.com/ypopivniak/msx/key"
)

// Item models the MSX Content Item Object. All exported fields have JSON tags
// and can be set directly on a zero value, but the fluent With* methods are
// the intended API.
//
// Pointer fields are used where an unset value differs from the Go zero
// value: default-true booleans (so explicit false round-trips), and numeric
// sentinels like -1 (so explicit 0 can round-trip without aliasing "unset").
type Item struct {
	ID   string   `json:"id,omitempty"`
	Type ItemType `json:"type,omitempty"`
	Key  key.Key  `json:"key,omitempty"`

	// omitzero (Go 1.24+) is required for struct-valued fields because
	// encoding/json's omitempty never treats a non-pointer struct as empty.
	Layout geom.Layout `json:"layout,omitzero"`
	Area   geom.Area   `json:"area,omitzero"`
	Offset geom.Offset `json:"offset,omitzero"`

	// Default-true flags — pointer so "WithDisplay(false)" survives marshal.
	Display  *bool `json:"display,omitempty"`
	Enable   *bool `json:"enable,omitempty"`
	Compress *bool `json:"compress,omitempty"`
	Shortcut *bool `json:"shortcut,omitempty"`
	// Default-false flags. Enumerate is a pointer because fixtures are seen
	// to write "enumerate": false explicitly even though it is the default.
	// The others stay as plain bool + omitempty until proven otherwise.
	Focus      bool  `json:"focus,omitempty"`
	Execute    bool  `json:"execute,omitempty"`
	Enumerate  *bool `json:"enumerate,omitempty"`
	Decompress bool  `json:"decompress,omitempty"`

	// Round is polymorphic (number|boolean) per the MSX spec. Use RoundInt or
	// RoundBool to construct; nil means "unset" and MSX applies its default.
	Round *RoundValue `json:"round,omitempty"`

	Break Break  `json:"break,omitempty"`
	Group string `json:"group,omitempty"`

	// Background is the default background image URL for the item.
	Background string `json:"background,omitempty"`

	Color       color.Color `json:"color,omitempty"`
	Title       string      `json:"title,omitempty"`
	TitleHeader string      `json:"titleHeader,omitempty"`
	TitleFooter string      `json:"titleFooter,omitempty"`
	Label       string      `json:"label,omitempty"`

	Icon     icon.Icon `json:"icon,omitempty"`
	IconSize icon.Size `json:"iconSize,omitempty"`

	Headline   string    `json:"headline,omitempty"`
	Text       textValue `json:"text,omitempty"`
	Alignment  Alignment `json:"alignment,omitempty"`
	Truncation string    `json:"truncation,omitempty"`
	Centration string    `json:"centration,omitempty"`
	Separation *float64  `json:"separation,omitempty"`

	Tag        string      `json:"tag,omitempty"`
	TagColor   color.Color `json:"tagColor,omitempty"`
	Badge      string      `json:"badge,omitempty"`
	BadgeColor color.Color `json:"badgeColor,omitempty"`
	Stamp      string      `json:"stamp,omitempty"`
	StampColor color.Color `json:"stampColor,omitempty"`

	Progress          *float64    `json:"progress,omitempty"`
	ProgressColor     color.Color `json:"progressColor,omitempty"`
	ProgressBackColor color.Color `json:"progressBackColor,omitempty"`
	WrapperColor      color.Color `json:"wrapperColor,omitempty"`

	Image             string      `json:"image,omitempty"`
	ImageFiller       ImageFiller `json:"imageFiller,omitempty"`
	ImageWidth        *float64    `json:"imageWidth,omitempty"`
	ImageHeight       *float64    `json:"imageHeight,omitempty"`
	ImageOverlay      *int        `json:"imageOverlay,omitempty"`
	ImagePreload      bool        `json:"imagePreload,omitempty"`
	ImageLabel        string      `json:"imageLabel,omitempty"`
	ImageColor        color.Color `json:"imageColor,omitempty"`
	ImageScreenFiller ImageFiller `json:"imageScreenFiller,omitempty"`
	ImageBoundary     bool        `json:"imageBoundary,omitempty"`

	PlayerLabel    string    `json:"playerLabel,omitempty"`
	ExtensionIcon  icon.Icon `json:"extensionIcon,omitempty"`
	ExtensionLabel string    `json:"extensionLabel,omitempty"`

	Action action.Action `json:"action,omitempty"`
	Data   any           `json:"data,omitempty"`
	// Properties values are string, number, or boolean (MSX accepts e.g.
	// "control:dim": true as well as "true"); Validate rejects other shapes.
	Properties map[string]any `json:"properties,omitempty"`

	// Live and Selection are strictly typed; Options is a polymorphic Panel.
	Live      *Live      `json:"live,omitempty"`
	Selection *Selection `json:"selection,omitempty"`
	Options   *Panel     `json:"options,omitempty"`
}

// Alignment is a free-form alignment expression (e.g. "left", "center",
// "title:center|text:justify"). Common constants are provided; any string is
// accepted to keep compound values (pipe-separated) expressible.
type Alignment string

const (
	AlignLeft    Alignment = "left"
	AlignCenter  Alignment = "center"
	AlignRight   Alignment = "right"
	AlignJustify Alignment = "justify"
)

// RoundValue models the polymorphic `round` property, which the MSX spec
// accepts as either a boolean (true/false) or an integer 0..3. We keep the
// original form so round-tripping a JSON document preserves whether the
// author wrote `true` or `1`, and so fixtures like `"round": false` can be
// reproduced byte-for-byte.
type RoundValue struct {
	isBool bool
	b      bool
	n      int
}

// RoundInt constructs a numeric round value. The valid range (0..3) is
// enforced by Item.Validate, not here.
func RoundInt(n int) *RoundValue { return &RoundValue{n: n} }

// RoundBool constructs a boolean round value.
func RoundBool(v bool) *RoundValue { return &RoundValue{isBool: true, b: v} }

// MarshalJSON emits either the bool or the integer depending on how the
// value was constructed.
func (r *RoundValue) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	if r.isBool {
		return json.Marshal(r.b)
	}
	return json.Marshal(r.n)
}

// UnmarshalJSON accepts either a JSON boolean or a JSON number; anything
// else is an error so round-trip tests catch unexpected shapes early.
func (r *RoundValue) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	// Try bool first — shorter payload, cheaper check.
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		r.isBool = true
		r.b = b
		return nil
	}
	var n int
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	r.isBool = false
	r.n = n
	return nil
}

// textValue stores the `text` property as a slice but marshals to a bare
// string when exactly one line is present (the common case), so output
// matches hand-written MSX content.
type textValue []string

func (t textValue) MarshalJSON() ([]byte, error) {
	switch len(t) {
	case 0:
		return []byte("null"), nil
	case 1:
		return json.Marshal(t[0])
	default:
		return json.Marshal([]string(t))
	}
}

// UnmarshalJSON accepts either a single JSON string or an array of strings,
// storing both forms as a []string internally. Round-trip preserves the
// original form because MarshalJSON above emits based on length.
func (t *textValue) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*t = nil
		return nil
	}
	// Array form: delegate to the default slice decoder.
	if data[0] == '[' {
		var arr []string
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		*t = arr
		return nil
	}
	// Scalar string form.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = []string{s}
	return nil
}

// idPattern is the validated character set for Item.ID per the wiki:
// letters, digits, "_", "-", ".".
var idPattern = regexp.MustCompile(`^[A-Za-z0-9_.\-]+$`)

// NewItem returns a fresh Item ready for chained configuration.
func NewItem() *Item { return &Item{} }

// ---- fluent setters (pointer receiver, returns *Item) ----

func (i *Item) WithID(id string) *Item                { i.ID = id; return i }
func (i *Item) WithType(t ItemType) *Item             { i.Type = t; return i }
func (i *Item) WithKey(k key.Key) *Item               { i.Key = k; return i }
func (i *Item) WithLayout(l geom.Layout) *Item        { i.Layout = l; return i }
func (i *Item) WithLayoutXYWH(x, y, w, h int) *Item   { i.Layout = geom.NewLayout(x, y, w, h); return i }
func (i *Item) WithArea(a geom.Area) *Item            { i.Area = a; return i }
func (i *Item) WithOffset(o geom.Offset) *Item        { i.Offset = o; return i }
func (i *Item) WithDisplay(v bool) *Item              { i.Display = &v; return i }
func (i *Item) WithEnable(v bool) *Item               { i.Enable = &v; return i }
func (i *Item) WithCompress(v bool) *Item             { i.Compress = &v; return i }
func (i *Item) WithShortcut(v bool) *Item             { i.Shortcut = &v; return i }
func (i *Item) WithFocus(v bool) *Item                { i.Focus = v; return i }
func (i *Item) WithExecute(v bool) *Item              { i.Execute = v; return i }
func (i *Item) WithEnumerate(v bool) *Item            { i.Enumerate = &v; return i }
func (i *Item) WithDecompress(v bool) *Item           { i.Decompress = v; return i }
func (i *Item) WithRound(v int) *Item                 { i.Round = RoundInt(v); return i }
func (i *Item) WithRoundBool(v bool) *Item            { i.Round = RoundBool(v); return i }
func (i *Item) WithBreak(b Break) *Item               { i.Break = b; return i }
func (i *Item) WithGroup(g string) *Item              { i.Group = g; return i }
func (i *Item) WithBackground(url string) *Item       { i.Background = url; return i }
func (i *Item) WithColor(c color.Color) *Item         { i.Color = c; return i }
func (i *Item) WithTitle(s string) *Item              { i.Title = s; return i }
func (i *Item) WithTitleHeader(s string) *Item        { i.TitleHeader = s; return i }
func (i *Item) WithTitleFooter(s string) *Item        { i.TitleFooter = s; return i }
func (i *Item) WithLabel(s string) *Item              { i.Label = s; return i }
func (i *Item) WithIcon(ic icon.Icon) *Item           { i.Icon = ic; return i }
func (i *Item) WithIconSize(s icon.Size) *Item        { i.IconSize = s; return i }
func (i *Item) WithHeadline(s string) *Item           { i.Headline = s; return i }
func (i *Item) WithText(lines ...string) *Item        { i.Text = lines; return i }
func (i *Item) WithAlignment(a Alignment) *Item       { i.Alignment = a; return i }
func (i *Item) WithTruncation(s string) *Item         { i.Truncation = s; return i }
func (i *Item) WithCentration(s string) *Item         { i.Centration = s; return i }
func (i *Item) WithSeparation(v float64) *Item        { i.Separation = &v; return i }
func (i *Item) WithTag(s string) *Item                { i.Tag = s; return i }
func (i *Item) WithTagColor(c color.Color) *Item      { i.TagColor = c; return i }
func (i *Item) WithBadge(s string) *Item              { i.Badge = s; return i }
func (i *Item) WithBadgeColor(c color.Color) *Item    { i.BadgeColor = c; return i }
func (i *Item) WithStamp(s string) *Item              { i.Stamp = s; return i }
func (i *Item) WithStampColor(c color.Color) *Item    { i.StampColor = c; return i }
func (i *Item) WithProgress(v float64) *Item          { i.Progress = &v; return i }
func (i *Item) WithProgressColor(c color.Color) *Item { i.ProgressColor = c; return i }
func (i *Item) WithProgressBack(c color.Color) *Item  { i.ProgressBackColor = c; return i }
func (i *Item) WithWrapperColor(c color.Color) *Item  { i.WrapperColor = c; return i }
func (i *Item) WithImage(url string) *Item            { i.Image = url; return i }
func (i *Item) WithImageFiller(f ImageFiller) *Item   { i.ImageFiller = f; return i }
func (i *Item) WithImageWidth(v float64) *Item        { i.ImageWidth = &v; return i }
func (i *Item) WithImageHeight(v float64) *Item       { i.ImageHeight = &v; return i }
func (i *Item) WithImageOverlay(v int) *Item          { i.ImageOverlay = &v; return i }
func (i *Item) WithImagePreload(v bool) *Item         { i.ImagePreload = v; return i }
func (i *Item) WithImageLabel(s string) *Item         { i.ImageLabel = s; return i }
func (i *Item) WithImageColor(c color.Color) *Item    { i.ImageColor = c; return i }
func (i *Item) WithImageScreen(f ImageFiller) *Item   { i.ImageScreenFiller = f; return i }
func (i *Item) WithImageBoundary(v bool) *Item        { i.ImageBoundary = v; return i }
func (i *Item) WithPlayerLabel(s string) *Item        { i.PlayerLabel = s; return i }
func (i *Item) WithExtensionIcon(ic icon.Icon) *Item  { i.ExtensionIcon = ic; return i }
func (i *Item) WithExtensionLabel(s string) *Item     { i.ExtensionLabel = s; return i }
func (i *Item) WithAction(a action.Action) *Item      { i.Action = a; return i }
func (i *Item) WithData(d any) *Item                  { i.Data = d; return i }
func (i *Item) WithProperties(p map[string]string) *Item {
	i.Properties = make(map[string]any, len(p))
	for k, v := range p {
		i.Properties[k] = v
	}
	return i
}
func (i *Item) WithProperty(k, v string) *Item { return i.WithPropertyValue(k, v) }

// WithPropertyValue sets a non-string extended property (bool or number).
func (i *Item) WithPropertyValue(k string, v any) *Item {
	if i.Properties == nil {
		i.Properties = map[string]any{}
	}
	i.Properties[k] = v
	return i
}
func (i *Item) WithOptions(p *Page) *Item        { i.Options = PanelFromPage(p); return i }
func (i *Item) WithOptionsPanel(p *Panel) *Item  { i.Options = p; return i }
func (i *Item) WithOptionsRoot(r *Content) *Item { i.Options = PanelFromContent(r); return i }
func (i *Item) WithLive(l *Live) *Item           { i.Live = l; return i }
func (i *Item) WithSelection(s *Selection) *Item { i.Selection = s; return i }

// Validate performs strict per-field checks. Cross-field rules that depend
// on the enclosing Page or Content (e.g. "layout is required unless the
// parent uses a template") are enforced in their respective Validate
// methods, not here.
func (i *Item) Validate() error {
	var errs []error
	if i.ID != "" && !idPattern.MatchString(i.ID) {
		errs = append(errs, validate.Newf("id", "contains invalid characters: %q", i.ID))
	}
	if i.Type != "" {
		errs = append(errs, validateOneOf("type", i.Type, allItemTypes))
	}
	if !i.Layout.IsZero() {
		errs = append(errs, validate.Prefix("layout", i.Layout.Validate()))
	}
	if !i.Area.IsZero() {
		errs = append(errs, validate.Prefix("area", i.Area.Validate()))
	}
	if !i.Offset.IsZero() {
		errs = append(errs, validate.Prefix("offset", i.Offset.Validate()))
	}
	if i.Round != nil && !i.Round.isBool && (i.Round.n < 0 || i.Round.n > 3) {
		errs = append(errs, validate.Newf("round", "numeric value must be in [0,3], got %d", i.Round.n))
	}
	if i.IconSize != "" {
		errs = append(errs, validateOneOf("iconSize", i.IconSize, icon.AllSizes))
	}
	if i.ImageFiller != "" {
		errs = append(errs, validateOneOf("imageFiller", i.ImageFiller, allFillers))
	}
	if i.ImageScreenFiller != "" {
		errs = append(errs, validateOneOf("imageScreenFiller", i.ImageScreenFiller, allFillers))
	}
	// Negative progress hides the bar but keeps it updatable (wiki).
	if i.Progress != nil && *i.Progress > 1 {
		errs = append(errs, validate.Newf("progress", "must be negative or in [0,1], got %g", *i.Progress))
	}
	// -1 is the documented default for separation/imageWidth/imageHeight.
	if i.Separation != nil && !isUnsetOrIn(*i.Separation, 1) {
		errs = append(errs, validate.Newf("separation", "must be -1 or in [0,1], got %g", *i.Separation))
	}
	if i.ImageOverlay != nil && (*i.ImageOverlay < -1 || *i.ImageOverlay > 4) {
		errs = append(errs, validate.Newf("imageOverlay", "must be in [-1,4], got %d", *i.ImageOverlay))
	}
	if i.ImageWidth != nil && !isUnsetOrIn(*i.ImageWidth, 12) {
		errs = append(errs, validate.Newf("imageWidth", "must be -1 or in [0,12], got %g", *i.ImageWidth))
	}
	if i.ImageHeight != nil && !isUnsetOrIn(*i.ImageHeight, 6) {
		errs = append(errs, validate.Newf("imageHeight", "must be -1 or in [0,6], got %g", *i.ImageHeight))
	}
	for k, v := range i.Properties {
		switch v.(type) {
		case string, bool, float64, float32, int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64, json.Number:
		default:
			errs = append(errs, validate.Newf("properties."+k, "must be a string, number, or boolean, got %T", v))
		}
	}
	if i.Live != nil {
		errs = append(errs, validate.Prefix("live", i.Live.Validate()))
	}
	if i.Selection != nil {
		errs = append(errs, validate.Prefix("selection", i.Selection.Validate()))
	}
	if i.Options != nil {
		errs = append(errs, validate.Prefix("options", i.Options.Validate()))
	}
	return validate.Join(errs...)
}

// isUnsetOrIn reports whether v is the -1 "unset" sentinel or in [0,max].
func isUnsetOrIn(v, max float64) bool {
	return v == -1 || (v >= 0 && v <= max)
}

package content

import (
	"strconv"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/geom"
	"github.com/ypopivniak/msx/internal/validate"
)

// Page models the MSX Content Page Object. It must contain at least one
// Item to be displayable (enforced by Validate).
type Page struct {
	// Default-true flags as pointers so explicit false is preserved.
	Display  *bool `json:"display,omitempty"`
	Compress *bool `json:"compress,omitempty"`

	// Pointer bools — allows round-tripping fixtures that explicitly set
	// the MSX default (e.g. "important": false) without losing the key.
	Important *bool `json:"important,omitempty"`
	Wrap      *bool `json:"wrap,omitempty"`

	// Transparent is 0|1|2; 0 is default so omitempty works.
	Transparent Tri `json:"transparent,omitempty"`

	Headline   string `json:"headline,omitempty"`
	Background string `json:"background,omitempty"`

	Area     geom.Area   `json:"area,omitzero"`
	Offset   geom.Offset `json:"offset,omitzero"`
	Position string      `json:"position,omitempty"`

	Template *Item   `json:"template,omitempty"`
	Items    []*Item `json:"items,omitempty"`

	Action action.Action `json:"action,omitempty"`
	Data   any           `json:"data,omitempty"`

	Options *Panel `json:"options,omitempty"`
	Caption string `json:"caption,omitempty"`

	// Sentinel-default (-1) fields kept as pointers to distinguish unset.
	CaptionUnderlay  *int     `json:"captionUnderlay,omitempty"`
	NavigationOffset *float64 `json:"navigationOffset,omitempty"`
	NavigationSpan   *float64 `json:"navigationSpan,omitempty"`
}

// NewPage returns a fresh Page ready for chained configuration.
func NewPage() *Page { return &Page{} }

func (p *Page) WithDisplay(v bool) *Page        { p.Display = &v; return p }
func (p *Page) WithCompress(v bool) *Page       { p.Compress = &v; return p }
func (p *Page) WithImportant(v bool) *Page      { p.Important = &v; return p }
func (p *Page) WithWrap(v bool) *Page           { p.Wrap = &v; return p }
func (p *Page) WithTransparent(t Tri) *Page     { p.Transparent = t; return p }
func (p *Page) WithHeadline(s string) *Page     { p.Headline = s; return p }
func (p *Page) WithBackground(url string) *Page { p.Background = url; return p }
func (p *Page) WithArea(a geom.Area) *Page      { p.Area = a; return p }
func (p *Page) WithOffset(o geom.Offset) *Page  { p.Offset = o; return p }
func (p *Page) WithPosition(s string) *Page     { p.Position = s; return p }
func (p *Page) WithTemplate(it *Item) *Page     { p.Template = it; return p }
func (p *Page) WithItems(items ...*Item) *Page  { p.Items = items; return p }

// AddItem appends an item and returns the page for chaining.
func (p *Page) AddItem(it *Item) *Page {
	p.Items = append(p.Items, it)
	return p
}

func (p *Page) WithAction(a action.Action) *Page { p.Action = a; return p }
func (p *Page) WithData(d any) *Page             { p.Data = d; return p }
func (p *Page) WithOptions(o *Page) *Page         { p.Options = PanelFromPage(o); return p }
func (p *Page) WithOptionsPanel(o *Panel) *Page   { p.Options = o; return p }
func (p *Page) WithOptionsRoot(r *Content) *Page  { p.Options = PanelFromContent(r); return p }
func (p *Page) WithCaption(s string) *Page        { p.Caption = s; return p }
func (p *Page) WithCaptionUnderlay(v int) *Page  { p.CaptionUnderlay = &v; return p }
func (p *Page) WithNavigationOffset(v float64) *Page {
	p.NavigationOffset = &v
	return p
}
func (p *Page) WithNavigationSpan(v float64) *Page {
	p.NavigationSpan = &v
	return p
}

// Validate enforces per-field rules plus the page-level invariant that
// Items must be non-empty (required by MSX for a page to render).
func (p *Page) Validate() error {
	var errs []error
	if len(p.Items) == 0 {
		errs = append(errs, validate.New("items", "page must have at least one item"))
	}
	if p.Transparent < 0 || p.Transparent > 2 {
		errs = append(errs, validate.Newf("transparent", "must be 0, 1, or 2, got %d", p.Transparent))
	}
	if p.CaptionUnderlay != nil {
		if v := *p.CaptionUnderlay; v < -1 || v > 1 {
			errs = append(errs, validate.Newf("captionUnderlay", "must be in [-1,1], got %d", v))
		}
	}
	if !p.Area.IsZero() {
		errs = append(errs, validate.Prefix("area", p.Area.Validate()))
	}
	if !p.Offset.IsZero() {
		errs = append(errs, validate.Prefix("offset", p.Offset.Validate()))
	}
	if p.Template != nil {
		errs = append(errs, validate.Prefix("template", p.Template.Validate()))
	}
	for idx, it := range p.Items {
		errs = append(errs, validate.Prefix(itemPath(idx), it.Validate()))
	}
	if p.Options != nil {
		errs = append(errs, validate.Prefix("options", p.Options.Validate()))
	}
	return validate.Join(errs...)
}

func itemPath(i int) string {
	return "items[" + strconv.Itoa(i) + "]"
}

package content

import (
	"strconv"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/internal/validate"
)

// Content is the MSX Content Root Object. Build it with New() and the
// fluent With* methods, then call Build() to validate and obtain a
// ready-to-marshal pointer.
type Content struct {
	Name      string `json:"name,omitempty"`
	Version   string `json:"version,omitempty"`
	Reference string `json:"reference,omitempty"`
	Flag      string `json:"flag,omitempty"`

	// Default-true flags — *bool so explicit false round-trips.
	Reuse   *bool `json:"reuse,omitempty"`
	Cache   *bool `json:"cache,omitempty"`
	Restore *bool `json:"restore,omitempty"`

	// Default-false flags.
	Important  bool `json:"important,omitempty"`
	Wrap       bool `json:"wrap,omitempty"`
	Compress   bool `json:"compress,omitempty"`
	Preselect  bool `json:"preselect,omitempty"`

	// Tri-state numbers with default 0.
	Refocus     Tri `json:"refocus,omitempty"`
	Transparent Tri `json:"transparent,omitempty"`

	Type    Type    `json:"type,omitempty"`
	Preload Preload `json:"preload,omitempty"`

	Headline   string `json:"headline,omitempty"`
	Background string `json:"background,omitempty"`
	Extension  string `json:"extension,omitempty"`
	Dictionary string `json:"dictionary,omitempty"`

	Template *Item   `json:"template,omitempty"`
	Items    []*Item `json:"items,omitempty"`
	Pages    []*Page `json:"pages,omitempty"`

	Header   *Page   `json:"header,omitempty"`
	Footer   *Page   `json:"footer,omitempty"`
	Inserts  []*Page `json:"inserts,omitempty"`
	Overlay  *Page   `json:"overlay,omitempty"`
	Underlay *Page   `json:"underlay,omitempty"`

	Action action.Action `json:"action,omitempty"`
	Data   any           `json:"data,omitempty"`
	Ready  *ReadyAction  `json:"ready,omitempty"`

	Options *Panel `json:"options,omitempty"`
	Caption string `json:"caption,omitempty"`

	CaptionUnderlay  *int     `json:"captionUnderlay,omitempty"`
	NavigationOffset *float64 `json:"navigationOffset,omitempty"`
	NavigationSpan   *float64 `json:"navigationSpan,omitempty"`
}

// ReadyAction is the small {action, data} envelope used by the `ready`
// property to trigger work once the content is fully displayed.
type ReadyAction struct {
	Action action.Action `json:"action,omitempty"`
	Data   any           `json:"data,omitempty"`
}

// New returns an empty Content ready for chained configuration.
func New() *Content { return &Content{} }

func (c *Content) WithName(s string) *Content      { c.Name = s; return c }
func (c *Content) WithVersion(s string) *Content   { c.Version = s; return c }
func (c *Content) WithReference(s string) *Content { c.Reference = s; return c }
func (c *Content) WithFlag(s string) *Content      { c.Flag = s; return c }
func (c *Content) WithReuse(v bool) *Content       { c.Reuse = &v; return c }
func (c *Content) WithCache(v bool) *Content       { c.Cache = &v; return c }
func (c *Content) WithRestore(v bool) *Content     { c.Restore = &v; return c }
func (c *Content) WithImportant(v bool) *Content   { c.Important = v; return c }
func (c *Content) WithWrap(v bool) *Content        { c.Wrap = v; return c }
func (c *Content) WithCompress(v bool) *Content    { c.Compress = v; return c }
func (c *Content) WithPreselect(v bool) *Content   { c.Preselect = v; return c }
func (c *Content) WithRefocus(t Tri) *Content      { c.Refocus = t; return c }
func (c *Content) WithTransparent(t Tri) *Content  { c.Transparent = t; return c }
func (c *Content) WithType(t Type) *Content        { c.Type = t; return c }
func (c *Content) WithPreload(p Preload) *Content  { c.Preload = p; return c }
func (c *Content) WithHeadline(s string) *Content  { c.Headline = s; return c }
func (c *Content) WithBackground(url string) *Content {
	c.Background = url
	return c
}
func (c *Content) WithExtension(s string) *Content  { c.Extension = s; return c }
func (c *Content) WithDictionary(s string) *Content { c.Dictionary = s; return c }
func (c *Content) WithTemplate(it *Item) *Content   { c.Template = it; return c }
func (c *Content) WithItems(items ...*Item) *Content {
	c.Items = items
	return c
}

// AddItem appends a template-driven item and returns the content for chaining.
func (c *Content) AddItem(it *Item) *Content { c.Items = append(c.Items, it); return c }

func (c *Content) WithPages(pages ...*Page) *Content { c.Pages = pages; return c }

// AddPage appends a page and returns the content for chaining.
func (c *Content) AddPage(p *Page) *Content { c.Pages = append(c.Pages, p); return c }

func (c *Content) WithHeader(p *Page) *Content     { c.Header = p; return c }
func (c *Content) WithFooter(p *Page) *Content     { c.Footer = p; return c }
func (c *Content) WithInserts(p ...*Page) *Content { c.Inserts = p; return c }
func (c *Content) WithOverlay(p *Page) *Content    { c.Overlay = p; return c }
func (c *Content) WithUnderlay(p *Page) *Content   { c.Underlay = p; return c }

func (c *Content) WithAction(a action.Action) *Content { c.Action = a; return c }
func (c *Content) WithData(d any) *Content             { c.Data = d; return c }
func (c *Content) WithReady(a action.Action, data any) *Content {
	c.Ready = &ReadyAction{Action: a, Data: data}
	return c
}

// WithOptions wraps a Page as the options panel. For a Content Root panel,
// use WithOptionsRoot instead.
func (c *Content) WithOptions(p *Page) *Content {
	c.Options = PanelFromPage(p)
	return c
}

// WithOptionsPanel accepts an already-built Panel (e.g. from the panel
// package) directly.
func (c *Content) WithOptionsPanel(p *Panel) *Content { c.Options = p; return c }

// WithOptionsRoot wraps a Content Root as the options panel.
func (c *Content) WithOptionsRoot(r *Content) *Content {
	c.Options = PanelFromContent(r)
	return c
}

func (c *Content) WithCaption(s string) *Content  { c.Caption = s; return c }
func (c *Content) WithCaptionUnderlay(v int) *Content {
	c.CaptionUnderlay = &v
	return c
}
func (c *Content) WithNavigationOffset(v float64) *Content {
	c.NavigationOffset = &v
	return c
}
func (c *Content) WithNavigationSpan(v float64) *Content {
	c.NavigationSpan = &v
	return c
}

// Validate performs strict checks on all fields, nested pages, items, and
// overlay/underlay/header/footer sub-pages. It enforces the template-vs-pages
// exclusivity rule: if Template is set, Pages must be empty and Items must
// be non-empty; otherwise, Pages must be non-empty.
func (c *Content) Validate() error {
	var errs []error

	if c.Type != "" {
		errs = append(errs, validateOneOf("type", c.Type, allTypes))
	}
	if c.Preload != "" {
		errs = append(errs, validateOneOf("preload", c.Preload, allPreloads))
	}
	if c.Refocus < 0 || c.Refocus > 2 {
		errs = append(errs, validate.Newf("refocus", "must be 0, 1, or 2, got %d", c.Refocus))
	}
	if c.Transparent < 0 || c.Transparent > 2 {
		errs = append(errs, validate.Newf("transparent", "must be 0, 1, or 2, got %d", c.Transparent))
	}
	if c.CaptionUnderlay != nil {
		if v := *c.CaptionUnderlay; v < -1 || v > 1 {
			errs = append(errs, validate.Newf("captionUnderlay", "must be in [-1,1], got %d", v))
		}
	}

	switch {
	case c.Template != nil && len(c.Pages) > 0:
		errs = append(errs, validate.New("", "template and pages are mutually exclusive"))
	case c.Template != nil && len(c.Items) == 0:
		errs = append(errs, validate.New("items", "required when template is set"))
	case c.Template == nil && len(c.Items) > 0:
		errs = append(errs, validate.New("items", "only valid when template is set; use pages instead"))
	}

	if c.Template != nil {
		errs = append(errs, validate.Prefix("template", c.Template.Validate()))
	}
	for i, it := range c.Items {
		errs = append(errs, validate.Prefix("items["+strconv.Itoa(i)+"]", it.Validate()))
	}
	for i, p := range c.Pages {
		errs = append(errs, validate.Prefix("pages["+strconv.Itoa(i)+"]", p.Validate()))
	}
	for i, p := range c.Inserts {
		errs = append(errs, validate.Prefix("inserts["+strconv.Itoa(i)+"]", p.Validate()))
	}
	if c.Header != nil {
		errs = append(errs, validate.Prefix("header", c.Header.Validate()))
	}
	if c.Footer != nil {
		errs = append(errs, validate.Prefix("footer", c.Footer.Validate()))
	}
	if c.Overlay != nil {
		errs = append(errs, validate.Prefix("overlay", c.Overlay.Validate()))
	}
	if c.Underlay != nil {
		errs = append(errs, validate.Prefix("underlay", c.Underlay.Validate()))
	}
	if c.Options != nil {
		errs = append(errs, validate.Prefix("options", c.Options.Validate()))
	}

	return validate.Join(errs...)
}

// Build validates the content strictly and returns it. The returned pointer
// is the same receiver; the method is a builder-style terminator that turns
// a construction chain into a (*Content, error) pair.
func (c *Content) Build() (*Content, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

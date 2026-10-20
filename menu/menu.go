package menu

import (
	"strconv"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/internal/validate"
)

// Menu models the MSX Menu Root Object — the top-level navigation document.
// A non-empty Menu slice is required; validation rejects an empty root.
type Menu struct {
	Name      string `json:"name,omitempty"`
	Version   string `json:"version,omitempty"`
	Reference string `json:"reference,omitempty"`
	Flag      string `json:"flag,omitempty"`

	// Default-true flags as pointers so explicit false round-trips.
	Reuse   *bool `json:"reuse,omitempty"`
	Cache   *bool `json:"cache,omitempty"`
	Restore *bool `json:"restore,omitempty"`

	// Tri-state numerics default to 0.
	Refocus     Tri `json:"refocus,omitempty"`
	Transparent Tri `json:"transparent,omitempty"`

	Style    Style    `json:"style,omitempty"`
	Logo     string   `json:"logo,omitempty"`
	LogoSize LogoSize `json:"logoSize,omitempty"`

	Headline   string `json:"headline,omitempty"`
	Background string `json:"background,omitempty"`
	Extension  string `json:"extension,omitempty"`
	Dictionary string `json:"dictionary,omitempty"`

	Menu []*Item `json:"menu,omitempty"`

	Action action.Action        `json:"action,omitempty"`
	Data   any                  `json:"data,omitempty"`
	Ready  *content.ReadyAction `json:"ready,omitempty"`

	Options *content.Panel `json:"options,omitempty"`
}

// New returns an empty Menu Root ready for chained configuration.
func New() *Menu { return &Menu{} }

func (m *Menu) WithName(s string) *Menu       { m.Name = s; return m }
func (m *Menu) WithVersion(s string) *Menu    { m.Version = s; return m }
func (m *Menu) WithReference(s string) *Menu  { m.Reference = s; return m }
func (m *Menu) WithFlag(s string) *Menu       { m.Flag = s; return m }
func (m *Menu) WithReuse(v bool) *Menu        { m.Reuse = &v; return m }
func (m *Menu) WithCache(v bool) *Menu        { m.Cache = &v; return m }
func (m *Menu) WithRestore(v bool) *Menu      { m.Restore = &v; return m }
func (m *Menu) WithRefocus(t Tri) *Menu       { m.Refocus = t; return m }
func (m *Menu) WithTransparent(t Tri) *Menu   { m.Transparent = t; return m }
func (m *Menu) WithStyle(s Style) *Menu       { m.Style = s; return m }
func (m *Menu) WithLogo(url string) *Menu     { m.Logo = url; return m }
func (m *Menu) WithLogoSize(s LogoSize) *Menu { m.LogoSize = s; return m }
func (m *Menu) WithHeadline(s string) *Menu   { m.Headline = s; return m }
func (m *Menu) WithBackground(url string) *Menu {
	m.Background = url
	return m
}
func (m *Menu) WithExtension(s string) *Menu  { m.Extension = s; return m }
func (m *Menu) WithDictionary(s string) *Menu { m.Dictionary = s; return m }
func (m *Menu) WithItems(items ...*Item) *Menu {
	m.Menu = items
	return m
}

// AddItem appends a menu item and returns the root for chaining.
func (m *Menu) AddItem(it *Item) *Menu { m.Menu = append(m.Menu, it); return m }

func (m *Menu) WithAction(a action.Action) *Menu { m.Action = a; return m }
func (m *Menu) WithData(d any) *Menu             { m.Data = d; return m }
func (m *Menu) WithReady(a action.Action, data any) *Menu {
	m.Ready = &content.ReadyAction{Action: a, Data: data}
	return m
}
func (m *Menu) WithOptions(p *content.Page) *Menu {
	m.Options = content.PanelFromPage(p)
	return m
}

// WithOptionsPanel accepts a pre-built Panel directly.
func (m *Menu) WithOptionsPanel(p *content.Panel) *Menu { m.Options = p; return m }

// WithOptionsRoot attaches a Content Root as the menu-level options panel.
func (m *Menu) WithOptionsRoot(r *content.Content) *Menu {
	m.Options = content.PanelFromContent(r)
	return m
}

// Validate performs strict per-field checks and recurses into every item
// and the options panel. The `menu` array must be non-empty.
func (m *Menu) Validate() error {
	var errs []error

	if len(m.Menu) == 0 {
		errs = append(errs, validate.New("menu", "must have at least one item"))
	}
	if m.Style != "" {
		errs = append(errs, validateOneOf("style", m.Style, allStyles))
	}
	if m.LogoSize != "" {
		errs = append(errs, validateOneOf("logoSize", m.LogoSize, allLogoSizes))
	}
	if m.Refocus < 0 || m.Refocus > 2 {
		errs = append(errs, validate.Newf("refocus", "must be 0, 1, or 2, got %d", m.Refocus))
	}
	if m.Transparent < 0 || m.Transparent > 2 {
		errs = append(errs, validate.Newf("transparent", "must be 0, 1, or 2, got %d", m.Transparent))
	}

	for i, it := range m.Menu {
		errs = append(errs, validate.Prefix("menu["+strconv.Itoa(i)+"]", it.Validate()))
	}
	if m.Options != nil {
		errs = append(errs, validate.Prefix("options", m.Options.Validate()))
	}

	return validate.Join(errs...)
}

// Build validates the menu strictly and returns it.
func (m *Menu) Build() (*Menu, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// validateOneOf mirrors the helper in the content package. Keeping a local
// alias avoids coupling menu to content's unexported helpers.
func validateOneOf[T comparable](path string, value T, allowed []T) error {
	return validate.OneOf(path, value, allowed...)
}

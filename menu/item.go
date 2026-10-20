package menu

import (
	"regexp"

	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/internal/validate"
)

// Item models the MSX Menu Item Object — one entry in the vertical menu list.
//
// Layout and presentation differ from a Content Item: menus use `label`,
// `icon`, optional `image`, plus an extension icon/label on the right.
// Selecting an item loads the corresponding content, either from `Data.Ref`
// or from `Data.Inline`.
type Item struct {
	ID   string   `json:"id,omitempty"`
	Type ItemType `json:"type,omitempty"`

	// Default-true flags as pointers so explicit false survives marshal.
	Display *bool `json:"display,omitempty"`
	Enable  *bool `json:"enable,omitempty"`

	Focus   bool `json:"focus,omitempty"`
	Execute bool `json:"execute,omitempty"`

	// Tri-state number (0|1|2) with default 0, so plain omitempty is correct.
	Transparent Tri `json:"transparent,omitempty"`

	Icon           icon.Icon `json:"icon,omitempty"`
	Image          string    `json:"image,omitempty"`
	Label          string    `json:"label,omitempty"`
	Background     string    `json:"background,omitempty"`
	ExtensionIcon  icon.Icon `json:"extensionIcon,omitempty"`
	ExtensionLabel string    `json:"extensionLabel,omitempty"`

	// LineColor applies only to separator-type items — validated below.
	LineColor color.Color `json:"lineColor,omitempty"`

	Data    *ItemData      `json:"data,omitempty"`
	Options *content.Panel `json:"options,omitempty"`
}

// idPattern is shared with content.Item validation — keep the regex in sync
// if the MSX spec ever changes.
var idPattern = regexp.MustCompile(`^[A-Za-z0-9_.\-]+$`)

// NewItem returns a fresh Item ready for chained configuration.
func NewItem() *Item { return &Item{} }

func (i *Item) WithID(id string) *Item               { i.ID = id; return i }
func (i *Item) WithType(t ItemType) *Item            { i.Type = t; return i }
func (i *Item) WithDisplay(v bool) *Item             { i.Display = &v; return i }
func (i *Item) WithEnable(v bool) *Item              { i.Enable = &v; return i }
func (i *Item) WithFocus(v bool) *Item               { i.Focus = v; return i }
func (i *Item) WithExecute(v bool) *Item             { i.Execute = v; return i }
func (i *Item) WithTransparent(t Tri) *Item          { i.Transparent = t; return i }
func (i *Item) WithIcon(ic icon.Icon) *Item          { i.Icon = ic; return i }
func (i *Item) WithImage(url string) *Item           { i.Image = url; return i }
func (i *Item) WithLabel(s string) *Item             { i.Label = s; return i }
func (i *Item) WithBackground(url string) *Item      { i.Background = url; return i }
func (i *Item) WithExtensionIcon(ic icon.Icon) *Item { i.ExtensionIcon = ic; return i }
func (i *Item) WithExtensionLabel(s string) *Item    { i.ExtensionLabel = s; return i }
func (i *Item) WithLineColor(c color.Color) *Item    { i.LineColor = c; return i }

// WithDataRef sets the menu item to load content from the given URL or
// request action (e.g. "content:https://…"). Clears any inline content.
func (i *Item) WithDataRef(ref string) *Item {
	i.Data = NewItemDataRef(ref)
	return i
}

// WithDataInline embeds a Content Root directly into the menu item.
func (i *Item) WithDataInline(c *content.Content) *Item {
	i.Data = NewItemDataInline(c)
	return i
}

// WithOptions attaches a Content Page as the menu-button options panel.
func (i *Item) WithOptions(p *content.Page) *Item {
	i.Options = content.PanelFromPage(p)
	return i
}

// WithOptionsPanel accepts a pre-built Panel (e.g. from the panel package).
func (i *Item) WithOptionsPanel(p *content.Panel) *Item { i.Options = p; return i }

// WithOptionsRoot attaches a Content Root as the menu-button options panel.
func (i *Item) WithOptionsRoot(r *content.Content) *Item {
	i.Options = content.PanelFromContent(r)
	return i
}

// Validate enforces per-field rules. Data is intentionally not required —
// MSX allows items that only trigger an options panel or are navigational
// stubs. Separator items may carry a lineColor; on other types it is not
// meaningful and is rejected.
func (i *Item) Validate() error {
	var errs []error
	if i.ID != "" && !idPattern.MatchString(i.ID) {
		errs = append(errs, validate.Newf("id", "contains invalid characters: %q", i.ID))
	}
	if i.Type != "" {
		errs = append(errs, validateOneOf("type", i.Type, allItemTypes))
	}
	if i.Transparent < 0 || i.Transparent > 2 {
		errs = append(errs, validate.Newf("transparent", "must be 0, 1, or 2, got %d", i.Transparent))
	}
	if i.LineColor != "" && i.Type != ItemTypeSeparator {
		errs = append(errs, validate.New("lineColor", "only valid on separator items"))
	}
	if i.Data != nil {
		errs = append(errs, validate.Prefix("data", i.Data.Validate()))
	}
	if i.Options != nil {
		errs = append(errs, validate.Prefix("options", i.Options.Validate()))
	}
	return validate.Join(errs...)
}

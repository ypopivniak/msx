package content

import (
	"encoding/json"

	"github.com/ypopivniak/msx/internal/validate"
)

// Panel wraps the value used by the `options` property, which MSX accepts
// as either a Content Page Object or a Content Root Object. Holding both
// shapes in one struct (with exactly one populated) preserves author
// intent on round-trip and keeps call sites strictly typed.
//
// Construction: use PanelFromPage or PanelFromContent; the dedicated
// builder methods on Content/Page/Item/Menu wrap these for you.
type Panel struct {
	page *Page
	root *Content
}

// PanelFromPage wraps a Content Page as an options-panel value.
func PanelFromPage(p *Page) *Panel {
	if p == nil {
		return nil
	}
	return &Panel{page: p}
}

// PanelFromContent wraps a Content Root as an options-panel value.
func PanelFromContent(c *Content) *Panel {
	if c == nil {
		return nil
	}
	return &Panel{root: c}
}

// Page returns the underlying Page, or nil if the panel holds a Root.
func (p *Panel) Page() *Page {
	if p == nil {
		return nil
	}
	return p.page
}

// Root returns the underlying Content, or nil if the panel holds a Page.
func (p *Panel) Root() *Content {
	if p == nil {
		return nil
	}
	return p.root
}

// MarshalJSON picks whichever side is populated.
func (p *Panel) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}
	if p.root != nil {
		return json.Marshal(p.root)
	}
	if p.page != nil {
		return json.Marshal(p.page)
	}
	return []byte("null"), nil
}

// rootOnlyKeys are Content Root properties that the Content Page Object does
// not have; any one of them marks an options payload as a Root.
var rootOnlyKeys = []string{
	"name", "version", "reference", "flag", "reuse", "cache", "restore",
	"preselect", "refocus", "type", "preload", "extension", "dictionary",
	"pages", "header", "footer", "inserts", "overlay", "underlay", "ready",
	"menu",
}

// UnmarshalJSON distinguishes Page vs Root by the presence of any top-level
// key unique to the Root Object (see rootOnlyKeys). A payload with none of
// those is treated as a Page.
func (p *Panel) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}
	isRoot := false
	for _, k := range rootOnlyKeys {
		if _, ok := probe[k]; ok {
			isRoot = true
			break
		}
	}
	if isRoot {
		var c Content
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		p.root = &c
		return nil
	}
	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return err
	}
	p.page = &page
	return nil
}

// Validate enforces that exactly one side is set and recurses into it.
func (p *Panel) Validate() error {
	if p == nil {
		return nil
	}
	switch {
	case p.page != nil && p.root != nil:
		return validate.New("", "panel must be either a page or a content root, not both")
	case p.page == nil && p.root == nil:
		return validate.New("", "panel must have either a page or a content root")
	case p.page != nil:
		return p.page.Validate()
	default:
		return p.root.Validate()
	}
}

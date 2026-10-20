// Package feed builds MSX URLs for feed-conversion services.
//
// MSX does not consume M3U/PLS/MRSS content directly. Instead it hits a
// conversion service (for M3U and PLS) or wraps an interaction-plugin
// URL (for MRSS) whose output it renders as a Content Root. This package
// exposes typed constructors for both flavors so callers don't assemble
// query strings or pipe-delimited plugin URLs by hand.
//
//   - M3U/PLS: https://msx.benzac.de/wiki/index.php?title=M3U/PLS_Files
//   - MRSS:    https://msx.benzac.de/wiki/index.php?title=MRSS_Feeds
package feed

import (
	"net/url"
	"strings"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/internal/validate"
)

// Default conversion-service endpoints, as documented on the wiki. Callers
// can override via the Endpoint field on Convert{M3U,PLS} if they host a
// mirror.
const (
	DefaultM3UEndpoint = "http://msx.benzac.de/services/m3u.php"
	DefaultPLSEndpoint = "http://msx.benzac.de/services/pls.php"
	// DefaultMRSSPlugin is the interaction-plugin URL MSX uses to render
	// MRSS feeds. It is appended after the MRSS feed URL with "@".
	DefaultMRSSPlugin = "http://msx.benzac.de/interaction/mrss.html"
)

// ContentKind selects the item category the M3U/PLS converter produces.
// "audio" is the service default.
type ContentKind string

const (
	KindVideo  ContentKind = "video"
	KindAudio  ContentKind = "audio"
	KindImage  ContentKind = "image"
	KindPlugin ContentKind = "plugin"
)

var allKinds = []ContentKind{KindVideo, KindAudio, KindImage, KindPlugin}

// Mode controls the converter's display-mode output.
type Mode string

const (
	ModeDefault        Mode = "default"
	ModeExtended       Mode = "extended"
	ModeExtendedHeader Mode = "extended-header"
)

var allModes = []Mode{ModeDefault, ModeExtended, ModeExtendedHeader}

// Context selects the contextual bucket the conversion populates.
type Context string

const (
	ContextDefault   Context = "default"
	ContextFavorites Context = "favorites"
)

var allContexts = []Context{ContextDefault, ContextFavorites}

// Convert models the parameters accepted by the M3U and PLS conversion
// services. Use NewM3U / NewPLS rather than zero-initializing.
type Convert struct {
	// Endpoint is the conversion-service URL. Defaults to
	// DefaultM3UEndpoint or DefaultPLSEndpoint depending on the
	// constructor you use.
	Endpoint string

	URL        string
	Kind       ContentKind
	Headline   string
	Background string
	Mode       Mode
	Context    Context
}

// NewM3U prepares a converter targeting the M3U service and the given
// playlist URL. All other parameters can be set via the returned struct.
func NewM3U(url string) *Convert {
	return &Convert{Endpoint: DefaultM3UEndpoint, URL: url}
}

// NewPLS prepares a converter targeting the PLS service.
func NewPLS(url string) *Convert {
	return &Convert{Endpoint: DefaultPLSEndpoint, URL: url}
}

func (c *Convert) WithKind(k ContentKind) *Convert    { c.Kind = k; return c }
func (c *Convert) WithHeadline(s string) *Convert     { c.Headline = s; return c }
func (c *Convert) WithBackground(url string) *Convert { c.Background = url; return c }
func (c *Convert) WithMode(m Mode) *Convert           { c.Mode = m; return c }
func (c *Convert) WithContext(ctx Context) *Convert   { c.Context = ctx; return c }
func (c *Convert) WithEndpoint(endpoint string) *Convert {
	c.Endpoint = endpoint
	return c
}

// Validate checks enum membership and that URL is present.
func (c *Convert) Validate() error {
	var errs []error
	if c.URL == "" {
		errs = append(errs, validate.New("url", "is required"))
	}
	if c.Endpoint == "" {
		errs = append(errs, validate.New("endpoint", "is required"))
	}
	if c.Kind != "" {
		errs = append(errs, validate.OneOf("type", c.Kind, allKinds...))
	}
	if c.Mode != "" {
		errs = append(errs, validate.OneOf("mode", c.Mode, allModes...))
	}
	if c.Context != "" {
		errs = append(errs, validate.OneOf("context", c.Context, allContexts...))
	}
	return validate.Join(errs...)
}

// URLString returns the fully assembled conversion URL with all set
// parameters query-encoded. Empty optional parameters are omitted so
// the resulting URL matches what a user would type by hand.
func (c *Convert) URLString() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	q := url.Values{}
	q.Set("url", c.URL)
	if c.Kind != "" {
		q.Set("type", string(c.Kind))
	}
	if c.Headline != "" {
		q.Set("headline", c.Headline)
	}
	if c.Background != "" {
		q.Set("background", c.Background)
	}
	if c.Mode != "" {
		q.Set("mode", string(c.Mode))
	}
	if c.Context != "" {
		q.Set("context", string(c.Context))
	}
	return c.Endpoint + "?" + q.Encode(), nil
}

// ContentAction returns the conversion URL wrapped as a content:<url>
// action string, ready to use as a menu-item or content-item data.
func (c *Convert) ContentAction() (action.Action, error) {
	s, err := c.URLString()
	if err != nil {
		return "", err
	}
	return action.Content(s), nil
}

// ---- MRSS ----

// MRSS assembles an MSX interaction-plugin URL for an MRSS feed. The
// wire format is "<feed>@<plugin>?<optional-pipe-params>" where the
// optional tail is a pipe-delimited list of banner/background/plugin/
// dictionary/brand strings in that fixed order.
type MRSS struct {
	FeedURL    string
	Plugin     string // defaults to DefaultMRSSPlugin
	Banner     string
	Background string
	PluginURL  string
	Dictionary string
	Brand      string
}

// NewMRSS prepares an MRSS wrapper for feedURL.
func NewMRSS(feedURL string) *MRSS {
	return &MRSS{FeedURL: feedURL, Plugin: DefaultMRSSPlugin}
}

func (m *MRSS) WithPlugin(p string) *MRSS       { m.Plugin = p; return m }
func (m *MRSS) WithBanner(url string) *MRSS     { m.Banner = url; return m }
func (m *MRSS) WithBackground(url string) *MRSS { m.Background = url; return m }
func (m *MRSS) WithPluginURL(url string) *MRSS  { m.PluginURL = url; return m }
func (m *MRSS) WithDictionary(url string) *MRSS { m.Dictionary = url; return m }
func (m *MRSS) WithBrand(color string) *MRSS    { m.Brand = color; return m }

// Validate ensures the feed URL is present.
func (m *MRSS) Validate() error {
	if m.FeedURL == "" {
		return validate.New("feedURL", "is required")
	}
	return nil
}

// URLString returns the MRSS composite URL. Trailing empty parameters
// are trimmed; internal empty slots are left as consecutive pipes per
// the MSX spec, so order is preserved when only later slots are set.
func (m *MRSS) URLString() (string, error) {
	if err := m.Validate(); err != nil {
		return "", err
	}
	plugin := m.Plugin
	if plugin == "" {
		plugin = DefaultMRSSPlugin
	}
	base := m.FeedURL + "@" + plugin

	tail := []string{m.Banner, m.Background, m.PluginURL, m.Dictionary, m.Brand}
	// Trim only trailing empties — internal empties must survive so the
	// positional meaning of later slots is preserved.
	last := len(tail)
	for last > 0 && tail[last-1] == "" {
		last--
	}
	if last == 0 {
		return base, nil
	}
	return base + "?" + strings.Join(tail[:last], "|"), nil
}

// InteractionAction wraps the composite URL as the MSX action
// "content:request:interaction:<composite>".
func (m *MRSS) InteractionAction() (action.Action, error) {
	s, err := m.URLString()
	if err != nil {
		return "", err
	}
	return action.Prefixed("content", "request:interaction:"+s), nil
}

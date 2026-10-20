// Package start models the MSX Start Object — the startup descriptor
// consumed by the MSX launcher (the typical file is named start.json and
// points the launcher at a menu or content URL).
//
// Unlike content and menu, Start has three genuinely required properties
// (name, version, parameter); Validate enforces their presence.
package start

import (
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/internal/validate"
)

// Welcome is the enum for the `welcome` property controlling the welcome
// screen shown at startup.
type Welcome string

const (
	WelcomeNone    Welcome = "none"
	WelcomeStart   Welcome = "start"
	WelcomeContent Welcome = "content"
)

var allWelcomes = []Welcome{WelcomeNone, WelcomeStart, WelcomeContent}

// Start is the MSX Start Object.
type Start struct {
	Name      string `json:"name,omitempty"`
	Version   string `json:"version,omitempty"`
	Parameter string `json:"parameter,omitempty"`

	// Welcome is optional — empty means "use MSX default (none)".
	Welcome Welcome `json:"welcome,omitempty"`

	// Launcher is the launcher-tile configuration. Nil omits the field.
	Launcher *Launcher `json:"launcher,omitempty"`
}

// Launcher configures how the MSX tile for this start parameter appears on
// the MSX home screen.
type Launcher struct {
	Icon  icon.Icon   `json:"icon,omitempty"`
	Image string      `json:"image,omitempty"`
	Color color.Color `json:"color,omitempty"`
}

// NewLauncher returns a fresh Launcher ready for chained configuration.
func NewLauncher() *Launcher { return &Launcher{} }

func (l *Launcher) WithIcon(ic icon.Icon) *Launcher   { l.Icon = ic; return l }
func (l *Launcher) WithImage(url string) *Launcher    { l.Image = url; return l }
func (l *Launcher) WithColor(c color.Color) *Launcher { l.Color = c; return l }

// New returns an empty Start ready for chained configuration.
func New() *Start { return &Start{} }

func (s *Start) WithName(v string) *Start        { s.Name = v; return s }
func (s *Start) WithVersion(v string) *Start     { s.Version = v; return s }
func (s *Start) WithParameter(v string) *Start   { s.Parameter = v; return s }
func (s *Start) WithWelcome(w Welcome) *Start    { s.Welcome = w; return s }
func (s *Start) WithLauncher(l *Launcher) *Start { s.Launcher = l; return s }

// WithLauncherIcon is a shortcut for the common case of wanting only an
// icon on the launcher tile; it allocates a Launcher if needed.
func (s *Start) WithLauncherIcon(ic icon.Icon) *Start {
	if s.Launcher == nil {
		s.Launcher = NewLauncher()
	}
	s.Launcher.Icon = ic
	return s
}

// Validate enforces required fields and enum membership.
func (s *Start) Validate() error {
	var errs []error
	if s.Name == "" {
		errs = append(errs, validate.New("name", "is required"))
	}
	if s.Version == "" {
		errs = append(errs, validate.New("version", "is required"))
	}
	if s.Parameter == "" {
		errs = append(errs, validate.New("parameter", "is required"))
	}
	if s.Welcome != "" {
		errs = append(errs, validate.OneOf("welcome", s.Welcome, allWelcomes...))
	}
	return validate.Join(errs...)
}

// Build validates and returns the Start object.
func (s *Start) Build() (*Start, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

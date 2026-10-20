package content

import (
	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/internal/validate"
)

// LiveType is the enum for `live.type` — the category of live behavior.
type LiveType string

const (
	LiveSchedule LiveType = "schedule"
	LiveLifetime LiveType = "lifetime"
	LiveAirtime  LiveType = "airtime"
	LivePlayback LiveType = "playback"
	LiveSetup    LiveType = "setup"
)

var allLiveTypes = []LiveType{
	LiveSchedule, LiveLifetime, LiveAirtime, LivePlayback, LiveSetup,
}

// DisplayOverride is the shared set of item-display properties that Live
// and its nested state objects can override when a live transition fires.
// It is embedded (flattened into the parent JSON object) rather than
// referenced as a sub-object to match the wiki's wire format.
type DisplayOverride struct {
	Title             string        `json:"title,omitempty"`
	TitleHeader       string        `json:"titleHeader,omitempty"`
	TitleFooter       string        `json:"titleFooter,omitempty"`
	Label             string        `json:"label,omitempty"`
	Icon              icon.Icon     `json:"icon,omitempty"`
	Color             color.Color   `json:"color,omitempty"`
	Headline          string        `json:"headline,omitempty"`
	Text              textValue     `json:"text,omitempty"`
	Tag               string        `json:"tag,omitempty"`
	TagColor          color.Color   `json:"tagColor,omitempty"`
	Badge             string        `json:"badge,omitempty"`
	BadgeColor        color.Color   `json:"badgeColor,omitempty"`
	Stamp             string        `json:"stamp,omitempty"`
	StampColor        color.Color   `json:"stampColor,omitempty"`
	Progress          *float64      `json:"progress,omitempty"`
	ProgressColor     color.Color   `json:"progressColor,omitempty"`
	ProgressBackColor color.Color   `json:"progressBackColor,omitempty"`
	WrapperColor      color.Color   `json:"wrapperColor,omitempty"`
	Image             string        `json:"image,omitempty"`
	ExtensionIcon     icon.Icon     `json:"extensionIcon,omitempty"`
	ExtensionLabel    string        `json:"extensionLabel,omitempty"`
	Action            action.Action `json:"action,omitempty"`
	Data              any           `json:"data,omitempty"`
}

// LiveState is a single live-lifecycle state (coming / over). It inherits
// the display-override fields and may carry its own execute envelope, which
// takes precedence over Live.Execute while the state is active.
type LiveState struct {
	DisplayOverride
	Execute *LiveExecute `json:"execute,omitempty"`
}

// NewLiveState returns an empty LiveState ready for chained configuration.
func NewLiveState() *LiveState { return &LiveState{} }

// LiveRunningState is the running-state object: a superset of LiveState
// that can schedule sub-transitions at quartile boundaries and at the
// ready/complete bookends.
type LiveRunningState struct {
	DisplayOverride
	Execute   *LiveExecute `json:"execute,omitempty"`
	Quartile1 *LiveState   `json:"quartile1,omitempty"`
	Quartile2 *LiveState   `json:"quartile2,omitempty"`
	Quartile3 *LiveState   `json:"quartile3,omitempty"`
	Quartile4 *LiveState   `json:"quartile4,omitempty"`
	Ready     *LiveState   `json:"ready,omitempty"`
	Complete  *LiveState   `json:"complete,omitempty"`
}

// NewLiveRunningState returns an empty LiveRunningState.
func NewLiveRunningState() *LiveRunningState { return &LiveRunningState{} }

// LiveExecute is the {action, data} envelope used under Live.Execute.
type LiveExecute struct {
	Action action.Action `json:"action,omitempty"`
	Data   any           `json:"data,omitempty"`
}

// Live models the MSX Live Object — periodic refresh, lifetime/airtime
// timing, and state-triggered property overrides.
//
// Most fields mirror the wiki one-for-one; display-override fields are
// embedded via DisplayOverride to avoid duplication across Live and its
// nested state types.
type Live struct {
	Type     LiveType `json:"type,omitempty"`
	From     int64    `json:"from,omitempty"`
	To       int64    `json:"to,omitempty"`
	Duration int64    `json:"duration,omitempty"`
	Delay    int64    `json:"delay,omitempty"`

	// Source uses the MSX "id"/"url"/"key"/"current"/"none" keyword set or
	// the extended "<kw>:<value>" / "resume:<kw>:<value>" forms. Kept as a
	// free-form string because the grammar is open-ended.
	Source string `json:"source,omitempty"`

	Coming  *LiveState        `json:"coming,omitempty"`
	Running *LiveRunningState `json:"running,omitempty"`
	Over    *LiveState        `json:"over,omitempty"`
	Execute *LiveExecute      `json:"execute,omitempty"`

	// Preload defaults to true; pointer lets "WithPreload(false)" survive.
	Preload *bool `json:"preload,omitempty"`

	// DisplayOverride provides the top-level property-override fields MSX
	// accepts directly on the Live object (outside any state sub-object).
	DisplayOverride
}

// NewLive returns an empty Live ready for chained configuration.
func NewLive() *Live { return &Live{} }

// ---- fluent setters for Live itself ----

func (l *Live) WithType(t LiveType) *Live              { l.Type = t; return l }
func (l *Live) WithFrom(ms int64) *Live                { l.From = ms; return l }
func (l *Live) WithTo(ms int64) *Live                  { l.To = ms; return l }
func (l *Live) WithDuration(ms int64) *Live            { l.Duration = ms; return l }
func (l *Live) WithDelay(ms int64) *Live               { l.Delay = ms; return l }
func (l *Live) WithSource(s string) *Live              { l.Source = s; return l }
func (l *Live) WithComing(s *LiveState) *Live          { l.Coming = s; return l }
func (l *Live) WithRunning(s *LiveRunningState) *Live  { l.Running = s; return l }
func (l *Live) WithOver(s *LiveState) *Live            { l.Over = s; return l }
func (l *Live) WithExecute(a action.Action, data any) *Live {
	l.Execute = &LiveExecute{Action: a, Data: data}
	return l
}
func (l *Live) WithPreload(v bool) *Live { l.Preload = &v; return l }

// Validate enforces enum membership (type), progress range for any state
// that sets it, and recurses into every state object.
func (l *Live) Validate() error {
	if l == nil {
		return nil
	}
	var errs []error
	if l.Type != "" {
		errs = append(errs, validateOneOf("type", l.Type, allLiveTypes))
	}
	if l.From < 0 {
		errs = append(errs, validate.Newf("from", "must be >= 0, got %d", l.From))
	}
	if l.To < 0 {
		errs = append(errs, validate.Newf("to", "must be >= 0, got %d", l.To))
	}
	if l.To != 0 && l.To < l.From {
		errs = append(errs, validate.Newf("to", "must be >= from (%d), got %d", l.From, l.To))
	}
	if l.Duration < 0 {
		errs = append(errs, validate.Newf("duration", "must be >= 0, got %d", l.Duration))
	}
	if l.Delay < 0 {
		errs = append(errs, validate.Newf("delay", "must be >= 0, got %d", l.Delay))
	}
	errs = append(errs, validate.Prefix("", l.DisplayOverride.validate()))
	if l.Coming != nil {
		errs = append(errs, validate.Prefix("coming", l.Coming.Validate()))
	}
	if l.Running != nil {
		errs = append(errs, validate.Prefix("running", l.Running.Validate()))
	}
	if l.Over != nil {
		errs = append(errs, validate.Prefix("over", l.Over.Validate()))
	}
	return validate.Join(errs...)
}

// Validate on a LiveState validates its embedded DisplayOverride.
func (s *LiveState) Validate() error {
	if s == nil {
		return nil
	}
	return s.DisplayOverride.validate()
}

// Validate on LiveRunningState recurses into each sub-state and its own
// display override.
func (r *LiveRunningState) Validate() error {
	if r == nil {
		return nil
	}
	errs := []error{r.DisplayOverride.validate()}
	type subState struct {
		path string
		s    *LiveState
	}
	for _, sub := range []subState{
		{"quartile1", r.Quartile1}, {"quartile2", r.Quartile2},
		{"quartile3", r.Quartile3}, {"quartile4", r.Quartile4},
		{"ready", r.Ready}, {"complete", r.Complete},
	} {
		if sub.s != nil {
			errs = append(errs, validate.Prefix(sub.path, sub.s.Validate()))
		}
	}
	return validate.Join(errs...)
}

// validate (lowercase) checks DisplayOverride's own constraints — primarily
// the progress range. It is used by the state Validate wrappers so errors
// carry the right field path.
func (d DisplayOverride) validate() error {
	var errs []error
	// Progress range mirrors Item: any negative value hides the bar (e.g. -1
	// "unset", -2 auto-update per Live docs); valid "manual" range is [0,1].
	if d.Progress != nil {
		v := *d.Progress
		if v > 1 {
			errs = append(errs, validate.Newf("progress", "must be negative or in [0,1], got %g", v))
		}
	}
	return validate.Join(errs...)
}

// ---- convenience setters on state types (mirror the display-field names) ----

func (s *LiveState) WithAction(a action.Action) *LiveState { s.Action = a; return s }
func (s *LiveState) WithData(d any) *LiveState             { s.Data = d; return s }
func (s *LiveState) WithHeadline(h string) *LiveState      { s.Headline = h; return s }
func (s *LiveState) WithTitle(t string) *LiveState         { s.Title = t; return s }
func (s *LiveState) WithExecute(a action.Action, data any) *LiveState {
	s.Execute = &LiveExecute{Action: a, Data: data}
	return s
}

func (r *LiveRunningState) WithAction(a action.Action) *LiveRunningState { r.Action = a; return r }
func (r *LiveRunningState) WithData(d any) *LiveRunningState             { r.Data = d; return r }
func (r *LiveRunningState) WithExecute(a action.Action, data any) *LiveRunningState {
	r.Execute = &LiveExecute{Action: a, Data: data}
	return r
}
func (r *LiveRunningState) WithQuartile1(s *LiveState) *LiveRunningState {
	r.Quartile1 = s
	return r
}
func (r *LiveRunningState) WithQuartile2(s *LiveState) *LiveRunningState {
	r.Quartile2 = s
	return r
}
func (r *LiveRunningState) WithQuartile3(s *LiveState) *LiveRunningState {
	r.Quartile3 = s
	return r
}
func (r *LiveRunningState) WithQuartile4(s *LiveState) *LiveRunningState {
	r.Quartile4 = s
	return r
}
func (r *LiveRunningState) WithReady(s *LiveState) *LiveRunningState {
	r.Ready = s
	return r
}
func (r *LiveRunningState) WithComplete(s *LiveState) *LiveRunningState {
	r.Complete = s
	return r
}

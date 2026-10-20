// Package dictionary models the MSX UI-translation dictionary file.
//
// A dictionary is a flat key -> (string | []string) map under a top-level
// `properties` object, with `name` and `version` metadata. See:
//   https://msx.benzac.de/wiki/index.php?title=Dictionary_Structure
//
// Values may be arrays of strings (since MSX 0.1.160) for multi-line
// entries; we expose Set for single-line and SetLines for array-valued
// keys and serialize each key in whichever form the caller chose.
package dictionary

import (
	"encoding/json"

	"github.com/ypopivniak/msx/internal/validate"
)

// Dictionary is a single translation file.
type Dictionary struct {
	Name       string `json:"name,omitempty"`
	Version    string `json:"version,omitempty"`
	Language   string `json:"language,omitempty"`
	Properties Props  `json:"properties,omitempty"`
}

// Props maps translation keys to their value (single string or []string).
// A dedicated type lets us customize marshaling so each key preserves
// the author's chosen shape without forcing every value into an array.
type Props map[string]Value

// Value is a single-line string or an array of strings. Exactly one side
// is populated; use StringValue/LinesValue to construct.
type Value struct {
	s     string
	lines []string
}

// StringValue constructs a single-line value.
func StringValue(s string) Value { return Value{s: s} }

// LinesValue constructs a multi-line (array) value.
func LinesValue(lines ...string) Value { return Value{lines: lines} }

// String returns the stored single-line value (empty if the value is
// multi-line; use Lines to distinguish).
func (v Value) String() string { return v.s }

// Lines returns the stored multi-line value (nil if single-line).
func (v Value) Lines() []string { return v.lines }

// MarshalJSON emits a bare string for single-line values and an array
// otherwise, matching the MSX wire format either way.
func (v Value) MarshalJSON() ([]byte, error) {
	if v.lines != nil {
		return json.Marshal(v.lines)
	}
	return json.Marshal(v.s)
}

// UnmarshalJSON accepts either a JSON string or an array of strings.
func (v *Value) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '[' {
		var arr []string
		if err := json.Unmarshal(data, &arr); err != nil {
			return err
		}
		v.lines = arr
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v.s = s
	return nil
}

// New returns an empty Dictionary ready for chained configuration.
func New() *Dictionary { return &Dictionary{Properties: Props{}} }

func (d *Dictionary) WithName(n string) *Dictionary    { d.Name = n; return d }
func (d *Dictionary) WithVersion(v string) *Dictionary { d.Version = v; return d }
func (d *Dictionary) WithLanguage(l string) *Dictionary {
	d.Language = l
	return d
}

// Set stores a single-line translation under key.
func (d *Dictionary) Set(key, value string) *Dictionary {
	if d.Properties == nil {
		d.Properties = Props{}
	}
	d.Properties[key] = StringValue(value)
	return d
}

// SetLines stores a multi-line translation (array form).
func (d *Dictionary) SetLines(key string, lines ...string) *Dictionary {
	if d.Properties == nil {
		d.Properties = Props{}
	}
	d.Properties[key] = LinesValue(lines...)
	return d
}

// Validate enforces the MSX requirement that name, version, and at least
// one property are present.
func (d *Dictionary) Validate() error {
	var errs []error
	if d.Name == "" {
		errs = append(errs, validate.New("name", "is required"))
	}
	if d.Version == "" {
		errs = append(errs, validate.New("version", "is required"))
	}
	if len(d.Properties) == 0 {
		errs = append(errs, validate.New("properties", "must contain at least one entry"))
	}
	return validate.Join(errs...)
}

// Build validates and returns the dictionary for marshaling.
func (d *Dictionary) Build() (*Dictionary, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return d, nil
}

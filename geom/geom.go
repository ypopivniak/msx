// Package geom models the MSX "x,y,w,h" string format used for layouts,
// areas, and positional offsets.
//
// MSX uses the same textual shape for three properties with slightly different
// value ranges and semantics:
//
//   - Layout: integer grid position, x in [0,11], y in [0,5], w in [1,12],
//     h in [1,6] (12x6 grid; compressed variants 8x6 are not modelled here).
//   - Area:   integer boundary area, same ranges as Layout but w/h may equal
//     the remaining space, and may be empty to mean "no area".
//   - Offset: fractional positional offset, x/w in [-11.0, 11.0] and
//     y/h in [-5.0, 5.0], used to nudge items relative to their layout.
//
// All three marshal to a string on the wire; their Go representation is a
// small struct so builder code stays readable and type-checked.
package geom

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ypopivniak/msx/internal/validate"
)

// Layout is an integer grid position/size. Use NewLayout to construct.
type Layout struct {
	X, Y, W, H int
}

// NewLayout builds a Layout value for a 12x6 grid. It does not validate;
// call Validate() or rely on the containing struct's Validate() to check.
func NewLayout(x, y, w, h int) Layout { return Layout{X: x, Y: y, W: w, H: h} }

func (l Layout) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", l.X, l.Y, l.W, l.H)
}

// Validate enforces the 12x6 grid ranges.
func (l Layout) Validate() error {
	return validate.Join(
		rangeInt("x", l.X, 0, 11),
		rangeInt("y", l.Y, 0, 5),
		rangeInt("w", l.W, 1, 12),
		rangeInt("h", l.H, 1, 6),
	)
}

// MarshalJSON emits the layout as its "x,y,w,h" string form.
func (l Layout) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(l.String())), nil
}

// UnmarshalJSON parses an "x,y,w,h" string back into the struct.
func (l *Layout) UnmarshalJSON(data []byte) error {
	x, y, w, h, err := parseIntXYWH(data)
	if err != nil {
		return err
	}
	l.X, l.Y, l.W, l.H = x, y, w, h
	return nil
}

// IsZero reports whether the layout is the zero value, so container
// structs can treat it as "unset" when marshaling with omitempty.
func (l Layout) IsZero() bool { return l == Layout{} }

// Area has the same ranges as Layout but models a boundary area; a zero value
// means "no area" and should be omitted from JSON output.
type Area struct {
	X, Y, W, H int
}

func NewArea(x, y, w, h int) Area { return Area{X: x, Y: y, W: w, H: h} }

func (a Area) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", a.X, a.Y, a.W, a.H)
}

func (a Area) Validate() error {
	return validate.Join(
		rangeInt("x", a.X, 0, 11),
		rangeInt("y", a.Y, 0, 5),
		rangeInt("w", a.W, 1, 12),
		rangeInt("h", a.H, 1, 6),
	)
}

func (a Area) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(a.String())), nil
}

func (a *Area) UnmarshalJSON(data []byte) error {
	x, y, w, h, err := parseIntXYWH(data)
	if err != nil {
		return err
	}
	a.X, a.Y, a.W, a.H = x, y, w, h
	return nil
}

func (a Area) IsZero() bool { return a == Area{} }

// Offset is a fractional position delta. Unlike Layout it allows negatives
// and non-integer values.
type Offset struct {
	X, Y, W, H float64
}

func NewOffset(x, y, w, h float64) Offset { return Offset{X: x, Y: y, W: w, H: h} }

func (o Offset) String() string {
	return fmt.Sprintf("%s,%s,%s,%s", f(o.X), f(o.Y), f(o.W), f(o.H))
}

func (o Offset) Validate() error {
	return validate.Join(
		rangeFloat("x", o.X, -11.0, 11.0),
		rangeFloat("y", o.Y, -5.0, 5.0),
		rangeFloat("w", o.W, -11.0, 11.0),
		rangeFloat("h", o.H, -5.0, 5.0),
	)
}

func (o Offset) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(o.String())), nil
}

func (o *Offset) UnmarshalJSON(data []byte) error {
	x, y, w, h, err := parseFloatXYWH(data)
	if err != nil {
		return err
	}
	o.X, o.Y, o.W, o.H = x, y, w, h
	return nil
}

func (o Offset) IsZero() bool { return o == Offset{} }

// f renders a float without trailing zeros so "1" serializes as "1" rather
// than "1.000000" — matching the form used across MSX examples.
func f(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func rangeInt(name string, v, lo, hi int) error {
	if v < lo || v > hi {
		return validate.Newf(name, "must be in [%d,%d], got %d", lo, hi, v)
	}
	return nil
}

func rangeFloat(name string, v, lo, hi float64) error {
	if v < lo || v > hi {
		return validate.Newf(name, "must be in [%g,%g], got %g", lo, hi, v)
	}
	return nil
}

// parseIntXYWH unmarshals data as a JSON string of the form "x,y,w,h" and
// returns the four integers. Shared by Layout and Area.
func parseIntXYWH(data []byte) (int, int, int, int, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return 0, 0, 0, 0, err
	}
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("geom: expected \"x,y,w,h\", got %q", s)
	}
	vals := [4]int{}
	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("geom: component %d of %q: %w", i, s, err)
		}
		vals[i] = n
	}
	return vals[0], vals[1], vals[2], vals[3], nil
}

// parseFloatXYWH is the Offset variant of parseIntXYWH.
func parseFloatXYWH(data []byte) (float64, float64, float64, float64, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return 0, 0, 0, 0, err
	}
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("geom: expected \"x,y,w,h\", got %q", s)
	}
	vals := [4]float64{}
	for i, p := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("geom: component %d of %q: %w", i, s, err)
		}
		vals[i] = v
	}
	return vals[0], vals[1], vals[2], vals[3], nil
}

// Package request models MSX HTTP request/response envelopes.
//
// MSX dynamic endpoints return a wrapped response object whose Status
// controls downstream behavior: 200 means MSX consumes Data as the actual
// payload (menu, content, etc.), anything else (still over HTTP 2xx, per
// the spec) surfaces Message to the user.
//
// This package also provides keyword-replacement helpers for URLs that
// embed MSX's {HOST}/{PLATFORM}/{LANGUAGE}/etc. placeholders, so callers
// can build parameterized URLs without hand-rolling strings.
package request

import (
	"fmt"
	"strings"

	"github.com/ypopivniak/msx/internal/validate"
)

// Response is the MSX HTTP response envelope.
//
// The status code is HTTP-compliant: 200 = success and Data carries the
// payload; 4xx/5xx = error and Message carries the description. Note that
// per the MSX spec, error responses are still delivered over an HTTP 2xx
// transport — the logical error status goes here in Status.
type Response struct {
	Status  int    `json:"status"`
	Text    string `json:"text,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// Success constructs a 200 OK response carrying the given payload.
func Success(data any) *Response {
	return &Response{Status: 200, Text: "OK", Data: data}
}

// Error constructs a non-2xx logical response with the given status code
// and error message. Status must be outside [200,299]; this is asserted
// at Validate time to catch accidental misuse.
func Error(status int, message string) *Response {
	return &Response{Status: status, Message: message}
}

// IsSuccess reports whether the response represents a 2xx status.
func (r *Response) IsSuccess() bool {
	return r != nil && r.Status >= 200 && r.Status < 300
}

// Validate enforces a sensible baseline: Status is set, Data appears only
// on success, Message appears only on errors. These are conventions from
// the wiki, not HTTP requirements — strict mode surfaces violations early.
func (r *Response) Validate() error {
	if r == nil {
		return validate.New("", "response is nil")
	}
	var errs []error
	if r.Status == 0 {
		errs = append(errs, validate.New("status", "is required"))
	} else if r.Status < 100 || r.Status > 599 {
		errs = append(errs, validate.Newf("status", "outside HTTP range [100,599], got %d", r.Status))
	}
	if r.IsSuccess() && r.Message != "" {
		errs = append(errs, validate.New("message", "must be empty on success responses"))
	}
	if !r.IsSuccess() && r.Data != nil {
		errs = append(errs, validate.New("data", "must be empty on error responses"))
	}
	return validate.Join(errs...)
}

// ---- URL keyword interpolation ----

// Keywords holds MSX URL keyword replacements. Any non-empty field is
// substituted into a template via Expand; empty fields are left as-is.
type Keywords struct {
	Host     string
	Secure   *bool
	Client   string
	Platform string
	ID       string
	Player   string
	Language string
	Width    int
	Height   int
}

// Expand replaces MSX keyword placeholders ({HOST}, {PLATFORM}, …) in url
// with the values from k. Unset fields are left as their literal
// placeholder so callers can layer further substitutions later.
func (k Keywords) Expand(url string) string {
	out := url
	if k.Host != "" {
		out = strings.ReplaceAll(out, "{HOST}", k.Host)
	}
	if k.Secure != nil {
		out = strings.ReplaceAll(out, "{SECURE}", boolStr(*k.Secure))
	}
	if k.Client != "" {
		out = strings.ReplaceAll(out, "{CLIENT}", k.Client)
	}
	if k.Platform != "" {
		out = strings.ReplaceAll(out, "{PLATFORM}", k.Platform)
	}
	if k.ID != "" {
		out = strings.ReplaceAll(out, "{ID}", k.ID)
	}
	if k.Player != "" {
		out = strings.ReplaceAll(out, "{PLAYER}", k.Player)
	}
	if k.Language != "" {
		out = strings.ReplaceAll(out, "{LANGUAGE}", k.Language)
	}
	if k.Width != 0 {
		out = strings.ReplaceAll(out, "{WIDTH}", fmt.Sprint(k.Width))
	}
	if k.Height != 0 {
		out = strings.ReplaceAll(out, "{HEIGHT}", fmt.Sprint(k.Height))
	}
	return out
}

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

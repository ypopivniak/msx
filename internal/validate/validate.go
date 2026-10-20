// Package validate provides shared validation primitives used across
// the msx packages. Errors carry a dotted field path so that callers
// can pinpoint exactly which property of a nested structure failed.
package validate

import (
	"errors"
	"fmt"
	"strings"
)

// FieldError is a single validation failure for a named field.
type FieldError struct {
	Path string
	Msg  string
}

func (e *FieldError) Error() string {
	if e.Path == "" {
		return e.Msg
	}
	return e.Path + ": " + e.Msg
}

// New returns a FieldError with a fixed message.
func New(path, msg string) error {
	return &FieldError{Path: path, Msg: msg}
}

// Newf returns a FieldError with a formatted message.
func Newf(path, format string, args ...any) error {
	return &FieldError{Path: path, Msg: fmt.Sprintf(format, args...)}
}

// Join merges non-nil errors into a single error using errors.Join.
// Nil entries are dropped; returns nil if all entries are nil.
func Join(errs ...error) error {
	filtered := errs[:0]
	for _, e := range errs {
		if e != nil {
			filtered = append(filtered, e)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return errors.Join(filtered...)
}

// Prefix rewrites the Path of every FieldError in err by prepending
// prefix + "." (or prefix alone if the error had no path). It is used
// to bubble nested-struct errors up through parent validators.
func Prefix(prefix string, err error) error {
	if err == nil {
		return nil
	}
	walk(err, func(fe *FieldError) {
		if fe.Path == "" {
			fe.Path = prefix
			return
		}
		fe.Path = prefix + "." + fe.Path
	})
	return err
}

func walk(err error, fn func(*FieldError)) {
	if err == nil {
		return
	}
	var fe *FieldError
	if errors.As(err, &fe) {
		fn(fe)
	}
	// errors.Join produces an error that exposes Unwrap() []error.
	type multi interface{ Unwrap() []error }
	if m, ok := err.(multi); ok {
		for _, e := range m.Unwrap() {
			walk(e, fn)
		}
	}
}

// OneOf checks that value is a member of allowed. Returns a FieldError
// listing the valid options otherwise.
func OneOf[T comparable](path string, value T, allowed ...T) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	parts := make([]string, len(allowed))
	for i, a := range allowed {
		parts[i] = fmt.Sprint(a)
	}
	return Newf(path, "value %v is not one of [%s]", value, strings.Join(parts, ", "))
}

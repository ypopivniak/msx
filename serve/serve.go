// Package serve provides thin HTTP helpers for exposing MSX documents.
//
// MSX's Setup Precondition requires every JSON endpoint to enable CORS so
// that the MSX web client can read it cross-origin. This package offers a
// CORS middleware with permissive defaults plus JSON handlers that marshal
// any Go value (and auto-call Validate if the value implements it) so
// callers don't repeat encoder boilerplate in every endpoint.
package serve

import (
	"encoding/json"
	"net/http"

	"github.com/ypopivniak/msx/request"
)

// Validator is the optional interface MSX builder types implement.
// JSON handlers check for it and call Validate before marshaling,
// turning any validation error into a 500 response so a misconfigured
// builder cannot accidentally serve an invalid document.
type Validator interface {
	Validate() error
}

// CORS wraps next with permissive CORS headers that satisfy the MSX
// setup precondition. It handles OPTIONS preflights with 204 and
// short-circuits before delegating to next for other verbs.
//
// Customize by composing your own middleware if a stricter origin
// policy is required; the defaults match MSX's documented minimum.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		h.Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// JSON serves v as application/json on every request, indenting with two
// spaces so the output matches the style MSX reference fixtures use.
// If v implements Validator and Validate fails, a 500 is returned with a
// plain-text error body — helpful during development, safe to leave on.
func JSON(v any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, v)
	}
}

// JSONFunc calls f on each request and serves its return value as JSON.
// Returned errors surface as 500 responses.
func JSONFunc(f func(r *http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := f(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, v)
	}
}

// Response serves f's return as a wrapped MSX Response envelope. This
// matches what MSX expects from dynamic (server-action) endpoints:
// always HTTP 200 on the wire, with the logical status embedded in the
// JSON body. A nil error produces a Success envelope; a returned
// *request.Response is served directly so callers can emit errors.
func Response(f func(r *http.Request) (*request.Response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := f(r)
		switch {
		case err != nil:
			resp = request.Error(500, err.Error())
		case resp == nil:
			resp = request.Success(nil)
		}
		if vErr := resp.Validate(); vErr != nil {
			http.Error(w, vErr.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, resp)
	}
}

// writeJSON is shared by all handlers — one place that decides on the
// status codes, the Content-Type, and the Validate short-circuit.
func writeJSON(w http.ResponseWriter, v any) {
	if vr, ok := v.(Validator); ok {
		if err := vr.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	body, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(body, '\n'))
}

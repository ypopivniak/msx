package serve_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/request"
	"github.com/ypopivniak/msx/serve"
)

func TestCORS_SetsHeadersAndPassesThrough(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "pong")
	})
	h := serve.CORS(mux)

	srv := httptest.NewServer(h)
	defer srv.Close()

	res, err := http.Get(srv.URL + "/ping")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("CORS origin missing, got %q", got)
	}
	body, _ := io.ReadAll(res.Body)
	if string(body) != "pong" {
		t.Errorf("body got %q, want pong", body)
	}
}

func TestCORS_PreflightReturns204(t *testing.T) {
	h := serve.CORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler must not run on OPTIONS")
	}))
	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("preflight must set Allow-Methods")
	}
}

// A builder object passed to JSON has its Validate() called automatically.
// A doc with a page that has no items is invalid; the handler must
// surface that as a 500 instead of serving a malformed document.
func TestJSON_AutoValidates(t *testing.T) {
	bad := content.New().AddPage(content.NewPage()) // page has no items
	rec := httptest.NewRecorder()
	serve.JSON(bad).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500 for invalid doc", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "at least one item") {
		t.Errorf("expected validation error body, got: %s", rec.Body.String())
	}
}

// A valid built object serves with 200 + the expected Content-Type.
func TestJSON_ServesValidDocument(t *testing.T) {
	c, err := content.New().
		AddPage(content.NewPage().AddItem(
			content.NewItem().WithTitle("x"),
		)).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	rec := httptest.NewRecorder()
	serve.JSON(c).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Errorf("Content-Type = %q", got)
	}
	var back map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &back); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if _, ok := back["pages"]; !ok {
		t.Errorf("expected pages in response: %s", rec.Body.String())
	}
}

func TestJSONFunc_CallsPerRequest(t *testing.T) {
	calls := 0
	h := serve.JSONFunc(func(_ *http.Request) (any, error) {
		calls++
		return map[string]int{"calls": calls}, nil
	})
	for i := 0; i < 3; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
	}
	if calls != 3 {
		t.Errorf("handler called %d times, want 3", calls)
	}
}

func TestJSONFunc_PropagatesErrors(t *testing.T) {
	h := serve.JSONFunc(func(_ *http.Request) (any, error) {
		return nil, io.ErrUnexpectedEOF
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

// Response handler wraps its payload in the MSX response envelope.
func TestResponse_WrapsSuccess(t *testing.T) {
	h := serve.Response(func(_ *http.Request) (*request.Response, error) {
		return request.Success(map[string]string{"hello": "world"}), nil
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d", rec.Code)
	}
	var env request.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Status != 200 {
		t.Errorf("env.Status = %d", env.Status)
	}
}

// A handler error collapses to an error envelope, still over HTTP 200
// per the MSX convention (logical status inside the body).
func TestResponse_HandlerErrorBecomesEnvelope(t *testing.T) {
	h := serve.Response(func(_ *http.Request) (*request.Response, error) {
		return nil, io.EOF
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 (error is embedded)", rec.Code)
	}
	var env request.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Status != 500 || env.Message == "" {
		t.Errorf("expected error envelope, got %+v", env)
	}
}

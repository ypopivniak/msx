package request_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/request"
)

func TestSuccess_Shape(t *testing.T) {
	r := request.Success(map[string]any{"headline": "Hi"})
	if !r.IsSuccess() {
		t.Fatal("Success should report IsSuccess")
	}
	if err := r.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, `"status":200`) {
		t.Errorf("status=200 missing: %s", s)
	}
	if !strings.Contains(s, `"text":"OK"`) {
		t.Errorf("text=OK missing: %s", s)
	}
	if !strings.Contains(s, `"data":{"headline":"Hi"}`) {
		t.Errorf("data missing: %s", s)
	}
}

func TestError_Shape(t *testing.T) {
	r := request.Error(404, "Not found")
	if r.IsSuccess() {
		t.Fatal("Error should not report success")
	}
	if err := r.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	b, _ := json.Marshal(r)
	s := string(b)
	if !strings.Contains(s, `"status":404`) {
		t.Errorf("status=404 missing: %s", s)
	}
	if !strings.Contains(s, `"message":"Not found"`) {
		t.Errorf("message missing: %s", s)
	}
	if strings.Contains(s, `"data"`) {
		t.Errorf("error must not include data: %s", s)
	}
}

func TestValidate_StatusMissing(t *testing.T) {
	r := &request.Response{}
	err := r.Validate()
	if err == nil || !strings.Contains(err.Error(), "status") {
		t.Errorf("expected status:is required, got %v", err)
	}
}

func TestValidate_StatusOutOfRange(t *testing.T) {
	r := &request.Response{Status: 42}
	err := r.Validate()
	if err == nil || !strings.Contains(err.Error(), "outside HTTP range") {
		t.Errorf("expected range error, got %v", err)
	}
}

func TestValidate_MessageOnSuccess(t *testing.T) {
	r := &request.Response{Status: 200, Message: "nope"}
	err := r.Validate()
	if err == nil || !strings.Contains(err.Error(), "message") {
		t.Errorf("expected message-on-success error, got %v", err)
	}
}

func TestValidate_DataOnError(t *testing.T) {
	r := &request.Response{Status: 500, Data: "payload"}
	err := r.Validate()
	if err == nil || !strings.Contains(err.Error(), "data") {
		t.Errorf("expected data-on-error error, got %v", err)
	}
}

func TestKeywords_Expand(t *testing.T) {
	tr := true
	k := request.Keywords{
		Host:     "msx.example",
		Secure:   &tr,
		Platform: "default",
		Language: "en",
		Width:    1920,
		Height:   1080,
	}
	tpl := "https://{HOST}/api?secure={SECURE}&p={PLATFORM}&l={LANGUAGE}&w={WIDTH}&h={HEIGHT}"
	got := k.Expand(tpl)
	want := "https://msx.example/api?secure=true&p=default&l=en&w=1920&h=1080"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestKeywords_LeavesUnsetAsIs(t *testing.T) {
	k := request.Keywords{Host: "msx.example"}
	got := k.Expand("{HOST}/{PLATFORM}")
	// Platform wasn't set, so {PLATFORM} remains.
	if got != "msx.example/{PLATFORM}" {
		t.Errorf("got %q", got)
	}
}

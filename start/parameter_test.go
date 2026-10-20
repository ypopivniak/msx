package start_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ypopivniak/msx/start"
)

func TestMenuParameter(t *testing.T) {
	if got := start.MenuParameter("https://x/y"); got != "menu:https://x/y" {
		t.Errorf("got %q", got)
	}
	if got := start.MenuUserParameter("https://x/y"); got != "menu:user:https://x/y" {
		t.Errorf("got %q", got)
	}
}

func TestContentParameter(t *testing.T) {
	if got := start.ContentParameter("https://x/y"); got != "content:https://x/y" {
		t.Errorf("got %q", got)
	}
	if got := start.ContentUserParameter("https://x/y"); got != "content:user:https://x/y" {
		t.Errorf("got %q", got)
	}
}

func TestRedirectParameter(t *testing.T) {
	if got := start.RedirectParameter("https://short/y"); got != "redirect:https://short/y" {
		t.Errorf("got %q", got)
	}
}

func TestIDParameter(t *testing.T) {
	if got := start.IDParameter("rgh", "benzac-de"); got != "id:rgh:benzac-de" {
		t.Errorf("got %q", got)
	}
}

// The parameter helpers compose with Start.WithParameter + Build and
// produce exactly the string MSX expects.
func TestStart_UsesParameterHelper(t *testing.T) {
	built, err := start.New().
		WithName("X").WithVersion("1").
		WithParameter(start.MenuParameter("https://benzac-de.github.io/msx/data/menu.json")).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	b, _ := json.Marshal(built)
	if !strings.Contains(string(b), `"parameter":"menu:https://benzac-de.github.io/msx/data/menu.json"`) {
		t.Errorf("parameter missing or wrong: %s", b)
	}
}

package expr_test

import (
	"testing"

	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/expr"
	"github.com/ypopivniak/msx/icon"
)

func TestText(t *testing.T) {
	if got := expr.Text(color.MSXWhite, "Hello"); got != "{txt:msx-white:Hello}" {
		t.Errorf("got %q", got)
	}
	// Empty color → uncolored form.
	if got := expr.Text("", "Hello"); got != "{txt:Hello}" {
		t.Errorf("got %q", got)
	}
}

func TestIcon(t *testing.T) {
	if got := expr.Icon(color.MSXWhite, icon.PlayArrow); got != "{ico:msx-white:play-arrow}" {
		t.Errorf("got %q", got)
	}
	if got := expr.Icon("", icon.Home); got != "{ico:home}" {
		t.Errorf("got %q", got)
	}
}

func TestCharAndColor(t *testing.T) {
	if got := expr.Char("copy"); got != "{chr:copy}" {
		t.Errorf("got %q", got)
	}
	if got := expr.CharColored(color.MSXRed, "copy"); got != "{chr:msx-red:copy}" {
		t.Errorf("got %q", got)
	}
	if got := expr.Col(color.MSXGlass); got != "{col:msx-glass}" {
		t.Errorf("got %q", got)
	}
}

func TestContextAndDict(t *testing.T) {
	if got := expr.Context("label"); got != "{context:label}" {
		t.Errorf("got %q", got)
	}
	if got := expr.Dict("welcome"); got != "{dict:welcome}" {
		t.Errorf("got %q", got)
	}
}

func TestNumberFamily(t *testing.T) {
	if got := expr.Number(1234); got != "{num:1234:format}" {
		t.Errorf("got %q", got)
	}
	if got := expr.NumberFormat(1234.5, ".,2"); got != "{num:1234.5:format:.,2}" {
		t.Errorf("got %q", got)
	}
	if got := expr.Date("now", "yyyy-MM-dd"); got != "{num:now:date:yyyy-MM-dd}" {
		t.Errorf("got %q", got)
	}
	if got := expr.Duration(5640000, "text:dhms"); got != "{num:5640000:duration:text:dhms}" {
		t.Errorf("got %q", got)
	}
}

func TestJoin(t *testing.T) {
	got := expr.Join(expr.Text(color.MSXWhite, "A"), expr.LineBreak, expr.Text(color.MSXRed, "B"))
	want := "{txt:msx-white:A}{br}{txt:msx-red:B}"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

// Mirror the welcome.json fixture string "{chr:ndash}" + {txt:msx-white:...}
// to verify the builder produces what the reference examples use.
func TestBuilder_MatchesFixtureIdiom(t *testing.T) {
	got := expr.New().
		Raw("Welcome to Media Station X ").
		Char("ndash").
		Raw(" GitHub Pages Example").
		String()
	want := "Welcome to Media Station X {chr:ndash} GitHub Pages Example"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestBuilder_ComplexRow(t *testing.T) {
	got := expr.New().
		Icon(color.MSXWhite, icon.PlayArrow).
		Raw(" For more information, see ").
		Text(color.MSXWhite, "https://example.com").
		Raw(".").
		Br().
		String()
	want := "{ico:msx-white:play-arrow} For more information, see {txt:msx-white:https://example.com}.{br}"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

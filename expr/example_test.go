package expr_test

import (
	"fmt"

	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/expr"
	"github.com/ypopivniak/msx/icon"
)

// Builder composes inline expressions for headlines, titles, and text.
func ExampleNew() {
	fmt.Println(expr.New().
		Icon(color.MSXWhite, icon.Movie).
		Raw(" Videos ").
		Text(color.MSXRed, "LIVE").
		String())
	// Output:
	// {ico:msx-white:movie} Videos {txt:msx-red:LIVE}
}

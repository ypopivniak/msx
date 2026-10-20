package action_test

import (
	"fmt"

	"github.com/ypopivniak/msx/action"
)

func ExampleVideo() {
	fmt.Println(action.Video("https://example.com/trailer.mp4"))
	// Output:
	// video:https://example.com/trailer.mp4
}

// Compound runs several actions in sequence.
func ExampleCompound() {
	fmt.Println(action.Compound(action.Info("Saved"), action.Back))
	// Output:
	// [info:Saved|back]
}

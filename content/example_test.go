package content_test

import (
	"encoding/json"
	"fmt"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/content"
)

// A template-driven content root: the template supplies shared properties
// and MSX lays the items out automatically.
func ExampleNew() {
	c, err := content.New().
		WithType(content.TypeList).
		WithHeadline("Videos").
		WithTemplate(content.NewItem().
			WithType(content.ItemTypeSeparate).
			WithLayoutXYWH(0, 0, 3, 3).
			WithColor(color.MSXGlass)).
		WithItems(
			content.NewItem().WithTitle("Clip 1").WithAction(action.Video("https://example.com/a.mp4")),
			content.NewItem().WithTitle("Clip 2").WithAction(action.Video("https://example.com/b.mp4")),
		).
		Build()
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println(string(b))
	// Output:
	// {
	//   "type": "list",
	//   "headline": "Videos",
	//   "template": {
	//     "type": "separate",
	//     "layout": "0,0,3,3",
	//     "color": "msx-glass"
	//   },
	//   "items": [
	//     {
	//       "title": "Clip 1",
	//       "action": "video:https://example.com/a.mp4"
	//     },
	//     {
	//       "title": "Clip 2",
	//       "action": "video:https://example.com/b.mp4"
	//     }
	//   ]
	// }
}

// Build validates the whole tree and reports each bad field with its path.
func ExampleContent_Build() {
	_, err := content.New().
		AddPage(content.NewPage().AddItem(content.NewItem().
			WithLayoutXYWH(0, 0, 4, 2).
			WithProgress(1.5))).
		Build()
	fmt.Println(err)
	// Output:
	// pages[0].items[0].progress: must be negative or in [0,1], got 1.5
}

// A live object switches item properties as a scheduled broadcast starts
// and ends.
func ExampleItem_WithLive() {
	it := content.NewItem().
		WithLayoutXYWH(0, 0, 6, 2).
		WithTitle("Evening News").
		WithLive(content.NewLive().
			WithType(content.LiveSchedule).
			WithFrom(1_767_200_000_000).
			WithTo(1_767_203_600_000).
			WithComing(content.NewLiveState().WithHeadline("Starts soon")).
			WithOver(content.NewLiveState().WithHeadline("Ended")))
	b, _ := json.Marshal(it)
	fmt.Println(string(b))
	// Output:
	// {"layout":"0,0,6,2","title":"Evening News","live":{"type":"schedule","from":1767200000000,"to":1767203600000,"coming":{"headline":"Starts soon"},"over":{"headline":"Ended"}}}
}

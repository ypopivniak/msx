package menu_test

import (
	"encoding/json"
	"fmt"

	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/menu"
)

// A menu item either references content by URL or embeds it inline.
func ExampleNew() {
	about := content.New().
		WithHeadline("About").
		AddPage(content.NewPage().AddItem(content.NewItem().
			WithLayoutXYWH(0, 0, 12, 2).
			WithText("Hello from MSX.")))

	m, err := menu.New().
		WithHeadline("My App").
		AddItem(menu.NewItem().
			WithIcon(icon.Movie).
			WithLabel("Videos").
			WithDataRef("https://example.com/msx/videos.json")).
		AddItem(menu.NewItem().
			WithIcon(icon.Info).
			WithLabel("About").
			WithDataInline(about)).
		Build()
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(b))
	// Output:
	// {
	//   "headline": "My App",
	//   "menu": [
	//     {
	//       "icon": "movie",
	//       "label": "Videos",
	//       "data": "https://example.com/msx/videos.json"
	//     },
	//     {
	//       "icon": "info",
	//       "label": "About",
	//       "data": {
	//         "headline": "About",
	//         "pages": [
	//           {
	//             "items": [
	//               {
	//                 "layout": "0,0,12,2",
	//                 "text": "Hello from MSX."
	//               }
	//             ]
	//           }
	//         ]
	//       }
	//     }
	//   ]
	// }
}

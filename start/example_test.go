package start_test

import (
	"encoding/json"
	"fmt"

	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/start"
)

// Build the start.json that MSX loads from http://{SERVER}/msx/start.json
// during setup. Name, version, and parameter are required.
func ExampleNew() {
	s, err := start.New().
		WithName("My App").
		WithVersion("1.0.0").
		WithParameter(start.MenuParameter("http://192.168.0.10:8080/msx/menu.json")).
		WithLauncherIcon(icon.Movie).
		Build()
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b))
	// Output:
	// {
	//   "name": "My App",
	//   "version": "1.0.0",
	//   "parameter": "menu:http://192.168.0.10:8080/msx/menu.json",
	//   "launcher": {
	//     "icon": "movie"
	//   }
	// }
}

// Build reports every missing required field at once.
func ExampleStart_Build() {
	_, err := start.New().WithName("My App").Build()
	fmt.Println(err)
	// Output:
	// version: is required
	// parameter: is required
}

func ExampleMenuParameter() {
	fmt.Println(start.MenuParameter("https://example.com/msx/menu.json"))
	fmt.Println(start.ContentParameter("https://example.com/msx/videos.json"))
	// Output:
	// menu:https://example.com/msx/menu.json
	// content:https://example.com/msx/videos.json
}

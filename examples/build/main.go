// Command build writes a static MSX deployment — start.json plus
// data/menu.json, data/videos.json, and data/dictionary.json — by composing
// the library's builder APIs.
//
// MSX's setup loads http://{SERVER}/msx/start.json and every other URL must
// be absolute, so write into a directory served at /msx/ and pass the
// matching -base:
//
//	go run ./examples/build -out ./site/msx -base http://192.168.0.10:8080/msx
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/dictionary"
	"github.com/ypopivniak/msx/expr"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/menu"
	"github.com/ypopivniak/msx/start"
)

func main() {
	out := flag.String("out", "./out", "directory to write JSON files into")
	base := flag.String("base", "https://example.invalid/msx", "base URL hosted by your server")
	flag.Parse()

	if err := os.MkdirAll(filepath.Join(*out, "data"), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	startDoc, err := buildStart(*base)
	if err != nil {
		log.Fatalf("start: %v", err)
	}
	menuDoc, err := buildMenu(*base)
	if err != nil {
		log.Fatalf("menu: %v", err)
	}
	videosDoc, err := buildVideos()
	if err != nil {
		log.Fatalf("videos: %v", err)
	}
	dictDoc, err := buildDictionary()
	if err != nil {
		log.Fatalf("dictionary: %v", err)
	}

	writeJSON(filepath.Join(*out, "start.json"), startDoc)
	writeJSON(filepath.Join(*out, "data", "menu.json"), menuDoc)
	writeJSON(filepath.Join(*out, "data", "videos.json"), videosDoc)
	writeJSON(filepath.Join(*out, "data", "dictionary.json"), dictDoc)

	fmt.Printf("wrote start.json, menu.json, videos.json, dictionary.json into %s\n", *out)
}

func buildStart(base string) (*start.Start, error) {
	return start.New().
		WithName("Example App").
		WithVersion("1.0.0").
		WithParameter(start.MenuParameter(base + "/data/menu.json")).
		WithLauncherIcon(icon.Movie).
		Build()
}

func buildMenu(base string) (*menu.Menu, error) {
	about := content.New().
		WithHeadline("About").
		AddPage(content.NewPage().AddItem(content.NewItem().
			WithType(content.ItemTypeDefault).
			WithLayoutXYWH(0, 0, 12, 2).
			WithHeadline("Example App").
			WithText("Built with github.com/ypopivniak/msx.")))

	return menu.New().
		WithHeadline("Example App").
		WithDictionary(base + "/data/dictionary.json").
		AddItem(menu.NewItem().
			WithType(menu.ItemTypeSeparator).
			WithLabel("Library")).
		AddItem(menu.NewItem().
			WithIcon(icon.Movie).
			WithLabel("Videos").
			WithDataRef(base + "/data/videos.json")).
		AddItem(menu.NewItem().
			WithIcon(icon.Info).
			WithLabel("About").
			WithDataInline(about)).
		Build()
}

func buildVideos() (*content.Content, error) {
	template := content.NewItem().
		WithType(content.ItemTypeSeparate).
		WithLayoutXYWH(0, 0, 3, 3).
		WithIcon(icon.Tinted(icon.TintWhiteSoft, string(icon.Movie))).
		WithColor(color.MSXGlass)

	row := func(n int, url string) *content.Item {
		title := fmt.Sprintf("Clip %d", n)
		return content.NewItem().
			WithTitle(title).
			WithPlayerLabel(title).
			WithAction(action.Video(url))
	}

	return content.New().
		WithType(content.TypePages).
		WithHeadline(expr.New().
			Icon(color.MSXWhite, icon.Movie).
			Raw(" Videos").
			String()).
		WithTemplate(template).
		WithItems(
			row(1, "https://example.invalid/media/a.mp4"),
			row(2, "https://example.invalid/media/b.mp4"),
			row(3, "https://example.invalid/media/c.mp4"),
		).
		Build()
}

func buildDictionary() (*dictionary.Dictionary, error) {
	return dictionary.New().
		WithName("English").
		WithVersion("1.0.0").
		WithLanguage("en").
		Set("caption:menu", "menu").
		Set("label:about", "About").
		SetLines("dictionary:message",
			"This service allows you to set or reset a custom dictionary.",
			" A dictionary from the startup data is ignored while a custom one is set.",
		).
		Build()
}

func writeJSON(path string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("marshal %s: %v", path, err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}

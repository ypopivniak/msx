# msx

[![Go Reference](https://pkg.go.dev/badge/github.com/ypopivniak/msx.svg)](https://pkg.go.dev/github.com/ypopivniak/msx)
[![CI](https://github.com/ypopivniak/msx/actions/workflows/ci.yml/badge.svg)](https://github.com/ypopivniak/msx/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/ypopivniak/msx)](https://github.com/ypopivniak/msx/releases)
[![License: MIT](https://img.shields.io/github/license/ypopivniak/msx)](LICENSE)

Strictly-typed Go API for producing and consuming [Media Station X](https://msx.benzac.de/) JSON documents. Builders for every MSX root object (Start, Menu, Content, Dictionary), polymorphic helpers for the tricky fields (Panel, Live, Selection, ItemData), and a small HTTP package for serving everything with the CORS headers MSX requires.

- Module path: `github.com/ypopivniak/msx`
- Go version: 1.24 (uses `omitzero` JSON tags and `errors.Join`)
- Every builder exposes `.Validate()` and `.Build() (*T, error)` — strict by default.

## Quickstart

### 1. See it on your TV

MSX runs on the TV (or in a browser) and pulls JSON from your server. Its setup asks for a server address and then loads `http://{SERVER}/msx/start.json`; that file points at a menu, and the menu points at content.

`examples/serve` builds those three documents in memory and serves them with the required CORS headers:

```sh
go run ./examples/serve -addr :8080
```

```
MSX listening on :8080
  start parameter: http://192.168.0.10:8080/msx/start.json
  in MSX: Settings → Start Parameter → Setup → enter 192.168.0.10:8080
```

1. Install **Media Station X** from your TV's app store (see [platform support](https://msx.benzac.de/info/?tab=PlatformSupport)).
2. Open **Settings → Start Parameter → Setup** and enter the address printed above (`192.168.0.10:8080`).
3. Confirm the start parameter. MSX loads the *Videos* menu, and the server logs each request it makes.

The server builds every link from your machine's LAN IP, because `localhost` on the TV means the TV itself. Pass `-host` to override it (e.g. `-host media.lan`).

### 2. Build documents in code

The documents are ordinary Go values. Each builder ends in `Build()`, which validates the whole tree and returns an error listing every bad field:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/ypopivniak/msx/action"
    "github.com/ypopivniak/msx/content"
    "github.com/ypopivniak/msx/icon"
    "github.com/ypopivniak/msx/menu"
    "github.com/ypopivniak/msx/start"
)

func main() {
    base := "http://192.168.0.10:8080/msx" // must be reachable from the TV

    startDoc, err := start.New().
        WithName("My App").
        WithVersion("1.0.0").
        WithParameter(start.MenuParameter(base + "/data/menu.json")).
        WithLauncherIcon(icon.Movie).
        Build()
    if err != nil {
        log.Fatal(err)
    }

    menuDoc, err := menu.New().
        WithHeadline("My App").
        AddItem(menu.NewItem().
            WithIcon(icon.Movie).
            WithLabel("Videos").
            WithDataRef(base + "/data/videos.json")).
        Build()
    if err != nil {
        log.Fatal(err)
    }

    videos, err := content.New().
        WithType(content.TypePages).
        WithHeadline("Videos").
        AddPage(content.NewPage().AddItem(
            content.NewItem().
                WithType(content.ItemTypeTeaser).
                WithLayoutXYWH(0, 0, 4, 6).
                WithTitle("Trailer").
                WithAction(action.Video("https://example.com/trailer.mp4")),
        )).
        Build()
    if err != nil {
        log.Fatal(err)
    }

    for _, doc := range []any{startDoc, menuDoc, videos} {
        b, _ := json.MarshalIndent(doc, "", "  ")
        fmt.Println(string(b))
    }
}
```

To serve them yourself, mount each one with `serve.JSON` at the path its URL points to, and wrap the mux in `serve.CORS`. [`examples/serve`](examples/serve/main.go) does exactly this.

### 3. Or write static files

`examples/build` writes the same kind of deployment to disk (plus an inline *About* page and a UI dictionary), for hosting on any web server:

```sh
go run ./examples/build -out ./site/msx -base http://192.168.0.10:8080/msx
```

Serve `./site` so that the files appear under `/msx/`. Any CORS-enabled server works, e.g. `npx http-server ./site --cors -p 8080`, or with this library:

```go
http.ListenAndServe(":8080", serve.CORS(http.FileServer(http.Dir("./site"))))
```

### Hosting requirements

- **CORS.** Every JSON response needs `Access-Control-Allow-Origin: *`, and `OPTIONS` preflights must succeed. `serve.CORS` handles both. Missing CORS headers are the most common reason MSX fails to load a file. See [Setup Precondition](https://msx.benzac.de/wiki/index.php?title=Setup_Precondition).
- **Absolute URLs.** Only `start.json` sits at a fixed path. Every URL inside your documents must be a full `http://` or `https://` URL that the TV can reach.
- **HTTPS.** The TV apps load plain HTTP. The browser version at `https://msx.benzac.de` only loads HTTPS files, so if you use it, set the security lock during setup and serve everything over HTTPS. See [Setup Start Parameter](https://msx.benzac.de/wiki/index.php?title=Setup_Start_Parameter).

## Package map

| Package | Purpose |
|---|---|
| `action` | MSX action strings — typed constructors (`Video`, `Content`, `Invalidate`, `Update`, `Focus`, `Player*`, …) and `Compound(a, b, c)` for `[a\|b\|c]`. |
| `color` | `Color` typed string + palette constants (`MSXGlass`, `MSXRed`, …). |
| `content` | Content Root, Page, Item with fluent builders; houses the polymorphic `Panel`, typed `Live`/`Selection`, `RoundValue`, and text string-or-array handling. |
| `dictionary` | UI-translation files; `Set(key, v)` / `SetLines(key, v...)` preserve single-vs-array shape on the wire. |
| `expr` | Inline expression builders — `Text`, `Icon`, `Char`, `Col`, `Context`, `Number`, `Date`, `Duration`, `Builder`. |
| `feed` | M3U / PLS conversion-service URL builder and MRSS interaction-plugin URL builder; `ContentAction()` wraps as `content:<url>`. |
| `geom` | `Layout` / `Area` / `Offset` structs; each marshals as `"x,y,w,h"` with range validation and `IsZero` for `omitzero`. |
| `icon` | Icon constants + `Tinted(tint, name)` for composite `"msx-white-soft:movie"` forms. |
| `key` | Typed constants for the `key` property (remote/keyboard bindings) + `Combine(keys...)` for `"a\|b"` groups. |
| `menu` | Menu Root and Menu Item; `ItemData` polymorphically holds a URL ref or an inline `content.Content`. |
| `panel` | Polymorphic Page/Root wrapper (alias for `content.Panel`) + `Info`/`Warn`/`Success`/`Error`/`Confirm`/`Message` dialog constructors. |
| `request` | MSX `Response` envelope (`Success`/`Error`) with validation, plus `Keywords{…}.Expand(url)` for `{HOST}` / `{PLATFORM}` / `{LANGUAGE}` interpolation. |
| `serve` | `CORS(next)` middleware plus `JSON`, `JSONFunc`, `Response` handlers that auto-validate anything implementing `Validator`. |
| `start` | Start Object (`name`/`version`/`parameter` required) with nested `Launcher`; `MenuParameter` / `ContentParameter` / `RedirectParameter` / `IDParameter` helpers. |

`internal/validate` holds the shared `FieldError`, `Join`, `Prefix`, and `OneOf` helpers used by every package's `Validate()`.

## Design conventions

**Builder pattern everywhere.** Each root type exposes `New()` → a zero-value pointer → fluent `With*` setters → `Build() (*T, error)` which runs `Validate()` and returns the same pointer. `Validate()` is always public if you prefer to separate construction from checking.

**Strict validation by default.** `Build()` returns an error listing every bad field with a dotted path (e.g. `pages[0].items[2].layout.x: must be in [0,11], got 99`). Multiple errors are joined via `errors.Join`.

**Optional fields that need three states use pointers.**

- Default-true booleans (`Display`, `Enable`, `Compress`, `Reuse`, `Cache`, `Restore`, …) are `*bool` so `WithX(false)` survives JSON round-trip through `omitempty`.
- Sentinel-default numerics (`CaptionUnderlay=-1`, `Progress=-1`, `Round` default `true/1`) are pointers so `0` doesn't alias "unset".

**Polymorphic JSON values use small wrapper types with custom `MarshalJSON`/`UnmarshalJSON`:**

- `content.RoundValue` — `RoundInt(n)` or `RoundBool(v)`.
- `content.Panel` — holds a `Page` **or** a `Content`; `UnmarshalJSON` dispatches on top-level keys.
- `menu.ItemData` — `Ref` URL string **or** inline `content.Content`.
- `content.textValue` (internal) — single string **or** string array, preserved on round-trip.
- `dictionary.Value` — same single/array story for translation entries.

**Struct-valued fields use `omitzero` (Go 1.24).** `geom.Layout`/`Area`/`Offset` are structs; `omitempty` never skips non-pointer structs, so the containing fields are tagged `omitzero` and the geom types implement `IsZero()`.

**Cross-field rules live in the enclosing `Validate`.** E.g. `content.Content.Validate` enforces the template-vs-pages exclusivity; `menu.Item.Validate` rejects `lineColor` on non-separator items.

**Panel is polymorphic (Page or Content Root).** `content.Page`/`Item`/`Content` and `menu.Menu`/`Item` each expose three setters for the `options` field: `WithOptions(*Page)`, `WithOptionsRoot(*Content)`, `WithOptionsPanel(*Panel)`.

## Examples

- [`examples/serve`](examples/serve/main.go) — build a start/menu/content set in memory and serve it under `/msx/` with CORS, using links built from the machine's LAN IP. This is the quickest way to see results on a TV.
- [`examples/build`](examples/build/main.go) — write `start.json`, `data/menu.json`, `data/videos.json`, and `data/dictionary.json` to disk for static hosting. It also shows inline menu content (`WithDataInline`) and a linked UI dictionary.

## Reference fixtures and round-trip tests

`testdata/reference/` mirrors the JSON files from [benzac-de/msx](https://github.com/benzac-de/msx): `start.json`, `data/menu.json`, `data/welcome.json`, the `data/base/*.json` set, and the `data/extended/*.json` pages. `TestContent_RoundTrip`, `TestMenu_RoundTrip`, and `TestStart_RoundTrip` unmarshal each fixture, `Validate`, re-marshal, and compare semantically (key order ignored; explicit-empty-string defaults normalized). Builder-reproduction tests (`TestBuilder_ReproducesStaticMenuFixture`, `TestMenu_BuilderReproduces_Reference`, `TestContent_BuilderReproduces_Videos`, `TestContent_BuilderReproduces_Welcome`, `TestStart_BuilderReproducesReference`) assert the fluent API can rebuild specific fixtures from scratch.

Run everything:

```sh
go test ./...
go vet  ./...
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the development workflow and how to add missing MSX fields. Security issues: see [SECURITY.md](SECURITY.md).

## License

MIT — see [LICENSE](LICENSE). The fixtures in `testdata/reference/` are copied from [benzac-de/msx](https://github.com/benzac-de/msx) and redistributed under the [Media Station X license](https://msx.benzac.de/info/?tab=License), not MIT; see [testdata/reference/README.md](testdata/reference/README.md).

This is an independent project and is not affiliated with or endorsed by Media Station X or its author.

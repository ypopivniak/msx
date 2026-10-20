// Command serve exposes a minimal MSX deployment over HTTP so a TV (or any
// MSX client on your network) can load it straight from this process — no
// files on disk required.
//
// MSX's setup asks for a server address and then loads
// http://{SERVER}/msx/start.json, so every document lives under /msx/.
// All other URLs must be absolute and reachable from the TV, which is why
// they are built from the machine's LAN address rather than "localhost".
//
//	go run ./examples/serve -addr :8080
//	# then in MSX: Settings → Start Parameter → Setup → enter <lan-ip>:8080
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ypopivniak/msx/action"
	"github.com/ypopivniak/msx/color"
	"github.com/ypopivniak/msx/content"
	"github.com/ypopivniak/msx/icon"
	"github.com/ypopivniak/msx/menu"
	"github.com/ypopivniak/msx/serve"
	"github.com/ypopivniak/msx/start"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	host := flag.String("host", "", "host the MSX client uses to reach this server (default: LAN IP)")
	flag.Parse()

	server, err := serverAddress(*addr, *host)
	if err != nil {
		log.Fatal(err)
	}
	h, err := newHandler("http://" + server + "/msx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MSX listening on %s\n", *addr)
	fmt.Printf("  start parameter: http://%s/msx/start.json\n", server)
	fmt.Printf("  in MSX: Settings → Start Parameter → Setup → enter %s\n", server)
	if err := http.ListenAndServe(*addr, logRequests(h)); err != nil {
		log.Fatal(err)
	}
}

// newHandler builds the documents with absolute URLs under base (e.g.
// "http://192.168.0.10:8080/msx") and mounts them at the matching paths.
func newHandler(base string) (http.Handler, error) {
	startDoc, err := start.New().
		WithName("Example App").
		WithVersion("1.0.0").
		WithParameter(start.MenuParameter(base + "/data/menu.json")).
		WithLauncherIcon(icon.Movie).
		Build()
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	menuDoc, err := menu.New().
		WithHeadline("Example App").
		AddItem(menu.NewItem().
			WithIcon(icon.Movie).
			WithLabel("Videos").
			WithDataRef(base + "/data/videos.json")).
		Build()
	if err != nil {
		return nil, fmt.Errorf("menu: %w", err)
	}

	videosDoc, err := content.New().
		WithType(content.TypePages).
		WithHeadline("Videos").
		WithTemplate(content.NewItem().
			WithType(content.ItemTypeSeparate).
			WithLayoutXYWH(0, 0, 3, 3).
			WithColor(color.MSXGlass)).
		WithItems(
			content.NewItem().
				WithTitle("Clip 1").
				WithAction(action.Video("https://example.invalid/a.mp4")),
			content.NewItem().
				WithTitle("Clip 2").
				WithAction(action.Video("https://example.invalid/b.mp4")),
		).
		Build()
	if err != nil {
		return nil, fmt.Errorf("videos: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/msx/start.json", serve.JSON(startDoc))
	mux.Handle("/msx/data/menu.json", serve.JSON(menuDoc))
	mux.Handle("/msx/data/videos.json", serve.JSON(videosDoc))
	return serve.CORS(mux), nil
}

// serverAddress returns the "host:port" the MSX client should use. The host
// is, in order: the -host flag, the host part of -addr if it names a
// specific interface, or the first non-loopback IPv4 address.
func serverAddress(addr, host string) (string, error) {
	h, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("invalid -addr %q: %w", addr, err)
	}
	if host == "" {
		if ip := net.ParseIP(h); h != "" && (ip == nil || !ip.IsUnspecified()) {
			host = h
		} else if host, err = lanIP(); err != nil {
			return "", err
		}
	}
	return net.JoinHostPort(host, port), nil
}

// lanIP prefers a private (RFC 1918) IPv4 address, falling back to any other
// non-loopback, non-link-local one.
func lanIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var fallback string
	for _, a := range addrs {
		n, ok := a.(*net.IPNet)
		if !ok || n.IP.To4() == nil || n.IP.IsLoopback() || n.IP.IsLinkLocalUnicast() {
			continue
		}
		if n.IP.IsPrivate() {
			return n.IP.String(), nil
		}
		if fallback == "" {
			fallback = n.IP.String()
		}
	}
	if fallback == "" {
		return "", fmt.Errorf("no LAN IPv4 address found; pass -host explicitly")
	}
	return fallback, nil
}

// logRequests prints each request so you can see the MSX client fetching
// start.json, the menu, and the content while setting it up.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

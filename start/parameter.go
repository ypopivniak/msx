package start

// Parameter string helpers.
//
// The `parameter` property of a Start Object is a simple prefixed-string
// grammar (not a structured object), so these helpers live here rather
// than in their own package. Each function returns a ready-to-use string
// suitable for WithParameter.

// MenuParameter returns "menu:<url>" — the common case for pointing the
// launcher at a Menu Root Object endpoint.
func MenuParameter(url string) string { return "menu:" + url }

// MenuUserParameter returns "menu:user:<url>" — loads the menu in the
// user-authenticated mode so MSX attaches stored credentials.
func MenuUserParameter(url string) string { return "menu:user:" + url }

// ContentParameter returns "content:<url>" — sets the launcher to open a
// Content Root Object directly, bypassing a menu.
func ContentParameter(url string) string { return "content:" + url }

// ContentUserParameter returns "content:user:<url>".
func ContentUserParameter(url string) string { return "content:user:" + url }

// RedirectParameter returns "redirect:<url>" — sends MSX to another URL
// that will in turn supply the real parameter. Used with URL shorteners.
func RedirectParameter(url string) string { return "redirect:" + url }

// IDParameter returns "id:<service>:<alias>" — the URL-shortener alias
// form (e.g. id:bly:example, id:rgh:benzac-de). The caller is responsible
// for knowing valid service identifiers; common ones include "bly"
// (Bitly), "cly" (Cuttly), and "rgh" (raw.githubusercontent).
func IDParameter(service, alias string) string { return "id:" + service + ":" + alias }

package serve_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/ypopivniak/msx/serve"
	"github.com/ypopivniak/msx/start"
)

// Mount documents at the paths their URLs point to and wrap the mux in
// CORS, which every MSX endpoint needs.
func ExampleJSON() {
	s, _ := start.New().
		WithName("My App").
		WithVersion("1.0.0").
		WithParameter(start.MenuParameter("http://192.168.0.10:8080/msx/menu.json")).
		Build()

	mux := http.NewServeMux()
	mux.Handle("/msx/start.json", serve.JSON(s))
	h := serve.CORS(mux)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/msx/start.json", nil))
	body, _ := io.ReadAll(rec.Body)
	fmt.Println(rec.Code, rec.Header().Get("Access-Control-Allow-Origin"))
	fmt.Print(string(body))
	// Output:
	// 200 *
	// {
	//   "name": "My App",
	//   "version": "1.0.0",
	//   "parameter": "menu:http://192.168.0.10:8080/msx/menu.json"
	// }
}

// Serve a directory of static JSON files (e.g. written by examples/build)
// with CORS enabled.
func ExampleCORS() {
	h := serve.CORS(http.FileServer(http.Dir("./site")))
	_ = h // http.ListenAndServe(":8080", h)
}

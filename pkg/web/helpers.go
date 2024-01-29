package web

import "net/http"

func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func htmxRedirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	w.Header().Set("HX-Location", url)
	w.WriteHeader(code)
}

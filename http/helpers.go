package http

import (
	"github.com/heyjorgedev/deploykit/http/view"
	"net/http"
)

func Error(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	view.ErrorPage(code).Render(w)
}

func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	if IsHTMX(r) {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, url, code)
}

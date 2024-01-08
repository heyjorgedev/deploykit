package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	handler := hello()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}

	if res.Body.String() != "Hello" {
		t.Errorf("expected body 'Hello', got %s", res.Body.String())
	}

}

func hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello"))
	}
}

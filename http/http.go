package http

import (
	"net/http"

	"github.com/jorgemurta/deploykit"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}

	eHandler, ok := err.(http.Handler)
	if ok {
		eHandler.ServeHTTP(w, r)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type ResponseWrapper[T any] struct {
	Data   T                          `json:"data"`
	Errors deploykit.ValidationErrors `json:"errors"`
}

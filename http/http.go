package http

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type Envelope[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`

	Data T `json:"data"`
}

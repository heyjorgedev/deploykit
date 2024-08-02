package http

import (
	"embed"
	"net/http"
)

//go:embed public/*
var publicFS embed.FS

func handlePublicDir() http.Handler {
	return http.FileServer(http.FS(publicFS))
}

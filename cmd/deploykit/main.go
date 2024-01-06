package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jorgemurta/deploykit/api"
)

func main() {
	srv := api.NewServer()

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = "0.0.0.0:9000"
	}

	err := http.ListenAndServe(addr, srv)
	if err != nil {
		log.Fatalln(err)
	}

}

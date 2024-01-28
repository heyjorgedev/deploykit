package main

import (
	"github.com/heyjorgedev/deploykit"
	"log"
)

func main() {
	app := deploykit.NewWithConfig()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

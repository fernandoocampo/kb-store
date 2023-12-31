package main

import (
	"log"
	"os"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/application"
)

func main() {
	app := application.NewServer()

	if err := app.Run(); err != nil {
		log.Printf("unable to start service: %s", err)
		os.Exit(-1)
	}

	log.Println("terminated")
}

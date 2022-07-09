package main

import (
	"log"
	"os"

	"github.com/kasulani/go-fundi/internal/app"
)

func main() {
	container := app.Container()

	if err := container.Invoke(app.Run); err != nil {
		log.Printf("failed to start application: %q\n", err)
		os.Exit(1)
	}

	defer container.Cleanup()
}

package main

import (
	"log"
	"os"

	"github.com/kasulani/go-fundi/internal/fundi"
)

func main() {
	container := fundi.Container("cli")

	if err := container.Invoke(fundi.StartCLI); err != nil {
		log.Printf("failed to start application: %q\n", err)
		os.Exit(1)
	}

	defer container.Cleanup()
}

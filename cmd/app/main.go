package main

import (
	"log"

	"github.com/v1adhope/time-tracker/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

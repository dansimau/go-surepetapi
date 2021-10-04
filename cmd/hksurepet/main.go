package main

import (
	"log"

	"github.com/dansimau/go-surepetapi/pkg/homekit"
)

func main() {
	hk, err := homekit.NewService()
	if err != nil {
		log.Fatal(err)
	}

	if err := hk.Run(); err != nil {
		log.Fatal(err)
	}
}

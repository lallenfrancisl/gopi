package main

import (
	"log"

	"github.com/lallenfrancisl/gopi/gopi"
)

func main() {
	api := gopi.New()

	route := api.Route("/users")
	route.
		Description("This is the description").
		Summary("This is the summary")

	js, err := api.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(js))
}

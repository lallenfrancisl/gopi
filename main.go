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

	type CreateUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	route.Get().
		Summary("List users").
		Tags([]string{"api"})

	route.Post().
		Summary("Create user").
		Tags([]string{"api"}).
		Body(&CreateUser{})

	js, err := api.MarshalJSONIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(js))
}

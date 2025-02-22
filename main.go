package main

import (
	"log"
	"net/http"
	"os"

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
		Body(&CreateUser{}).
		Response(http.StatusOK, &CreateUser{})

	js, err := api.MarshalJSONIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("schema.json", js, os.FileMode(os.O_RDWR))
	if err != nil {
		log.Fatal(err)
	}
}

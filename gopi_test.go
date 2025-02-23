package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lallenfrancisl/gopi/gopi"
)

func main() {
	api := gopi.New()

	route := api.Route("/users")
	route.
		Description("This is the description").
		Summary("This is the summary")

	type Address struct {
		state   string
		country string
	}

	type CreateUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Time struct {
		Hour    int
		Minutes int
		Seconds int
	}

	type User struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt Time      `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Address   Address
	}

	route.Get().
		Summary("List users").
		Tags([]string{"api"}).
		Response(http.StatusOK, time.Time{})

	route.Post().
		Summary("Create user").
		Tags([]string{"api"}).
		Body(&CreateUser{}).
		Response(http.StatusOK, 1.5).
		Response(http.StatusBadRequest, &User{})

	js, err := api.MarshalJSONIndent("", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./build/schema.json", js, os.FileMode(os.O_RDWR))
	if err != nil {
		log.Fatal(err)
	}
}

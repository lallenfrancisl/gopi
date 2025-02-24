package gopi_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/lallenfrancisl/gopi"
)

func TestGopi(t *testing.T) {
	api := gopi.New()

	route := api.Route("/users")
	route.
		Description("This is the description").
		Summary("This is the summary")

	type Address struct {
		State   string
		Country string
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
		Address   Address   `json:"address"`
	}

	type envelope map[string]any

	route.Get().
		Summary("List users").
		Tags([]string{"api"}).
		Response(http.StatusOK, envelope{"user": &User{}})

	route.Post().
		Summary("Create user").
		Tags([]string{"api"}).
		Body(&CreateUser{}).
		Response(http.StatusOK, envelope{"user": &CreateUser{}}).
		Response(http.StatusBadRequest, envelope{"error": ""}).
		Response(http.StatusInternalServerError, map[string]string{})

	js, err := api.MarshalJSONIndent("", "    ")
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("./build/schema.json", js, os.FileMode(os.O_TRUNC))
	if err != nil {
		t.Fatal(err)
	}
}

# gopi (In active development)

A fluent api for building OpenAPI 3.0 schema intuitively. Still in active development, so expect breaking api changes

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/lallenfrancisl/gopi"
)

var docs *gopi.Gopi = gopi.New()

func main () {
	docs.
		Title("Greenlight movie database RESTful API").
		Description(`
			Greenlight is an api for a service like IMDB, where users can
			add, list and edit details about movies. I built this to learn building
			web APIs in Go. The api OpenAPI API definition of this was created using
			https://github.com/lallenfrancisl/gopi, a tool that I made. And the documentation
			UI is rendered using https://scalar.com
		`).
		Contact(gopi.ContactDef{
			Name: "Allen Francis",
		}).
		Version("1.0.0").
		DefineTag(gopi.TagDef{
			Name:        "Movies",
			Description: "APIs for managing movies",
		})

    AddDocs()

	js, err := docs.MarshalJSONIndent("", "    ")
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	err = os.WriteFile("./docs/swagger.json", js, os.FileMode(os.O_TRUNC))
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}

func AddDocs() {
	docs.Route("/v1/movies").Post().
		Summary("Create a new movie").
		Tags([]string{"Movies"}).
		Body(&createMoviePayload{}).
		Response(
			http.StatusOK,
			envelope{"movie": &data.Movie{}},
		)

	docs.Route("/v1/movies/{id}").Get().
		Summary("Get a movie by id").
		Params(
			gopi.PathParam("id", 0).
				Description("Id of the movie").
				Required(),
		).
		Tags([]string{"Movies"}).
		Response(http.StatusOK, envelope{"movie": &data.Movie{}})

	docs.Route("/v1/movies/{id}").Patch().
		Summary("Update a movie by id").
		Tags([]string{"Movies"}).
		Params(
			gopi.PathParam("id", 0).
				Description("Id of the movie").
				Required(),
		).
		Body(updateMoviePayload{}).
		Response(http.StatusOK, envelope{"movie": data.Movie{}})

	docs.Route("/v1/movies/{id}").Delete().
		Summary("Delete a movie by id").
		Tags([]string{"Movies"}).
		Params(
			gopi.PathParam("id", 0).
				Description("Id of the movie").
				Required(),
		).
		Response(http.StatusOK, envelope{"message": ""})

	docs.Route("/v1/movies").Get().
		Summary("List all the movies").
		Tags([]string{"Movies"}).
		Params(
			gopi.QueryParam("title", "").
				Description("Search by title"),
			gopi.QueryParam("genres", []string{}).
				Description("Filter by list of genres"),
			gopi.QueryParam("page", 0).
				Description("Page number"),
			gopi.QueryParam("page_size", 0).
				Description("Number of items in each page"),
			gopi.QueryParam("sort", "").
				Description("Sort by given field name and direction"),
		).
		Response(
			http.StatusOK,
			envelope{"movies": []data.Movie{}, "metadata": data.Metadata{}},
		)
}
```

## Demo

A live demo of this used along with [Scalar](https://scalar.com) can be seen [here](https://api.allenfrancis.me/greenlight/v1/docs).

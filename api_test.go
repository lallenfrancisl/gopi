package main

import (
	"os"
	"testing"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func TestLibOpenAPI(t *testing.T) {
	sample, err := os.ReadFile("sample.yaml")
	if err != nil {
		t.Error(err)
	}

	doc, err := libopenapi.NewDocument(sample)
	if err != nil {
		t.Error(err)
	}

	dom, errs := doc.BuildV3Model()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Log(err)
		}
	}

	t.Log(dom.Model.Paths.PathItems)
}

func TestRouteModel(t *testing.T) {
	api := NewApi("test api", "1.0.0")

	doc := api.Route(
		CreateRouteParams{
			path:        "/users",
			description: "test description",
			summary:     "test summary",
		},
	)

	doc.Get(Operation{
		Operation: v3.Operation{
			Description: "List users",
			Summary:     "List all the registered user",
		},
	})

	js, err := api.dom.RenderJSON("    ")
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(string(js))
}

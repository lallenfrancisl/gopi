package main

import (
	"os"
	"testing"

	"github.com/pb33f/libopenapi"
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
	t.Log(api)
}

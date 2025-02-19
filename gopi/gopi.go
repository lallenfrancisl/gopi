package gopi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func newSpec(name string) *openapi3.T {
	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:      name,
			Version:    "0.0.0",
			Extensions: map[string]interface{}{},
		},
		Components: &openapi3.Components{
			Schemas:    make(openapi3.Schemas),
			Extensions: map[string]interface{}{},
		},
		Paths:      &openapi3.Paths{},
		Extensions: map[string]interface{}{},
	}
}

type Gopi struct {
	spec *openapi3.T
}

// Create a documentation object for a route
func (gopi *Gopi) Route(path string) *Route {
	pathItem := gopi.spec.Paths.Find(path)

	if pathItem == nil {
		pathItem = &openapi3.PathItem{}
		gopi.spec.Paths.Set(path, pathItem)
	}

	return &Route{
		gopi: gopi,
		Path: path,
	}
}

func (gopi *Gopi) MarshalJSON() ([]byte, error) {
	return gopi.spec.MarshalJSON()
}

type Route struct {
	gopi *Gopi
	Path string
}

// Set the summary of a route
func (route *Route) Summary(text string) *Route {
	pathItem := route.gopi.spec.Paths.Find(route.Path)

	if pathItem == nil {
		return route
	}

	pathItem.Summary = text

	return route
}

// Set the description of a route
func (route *Route) Description(text string) *Route {
	pathItem := route.gopi.spec.Paths.Find(route.Path)

	if pathItem == nil {
		return route
	}

	pathItem.Description = text

	return route
}

func New() *Gopi {
	spec := newSpec("test")

	api := &Gopi{
		spec: spec,
	}

	return api
}

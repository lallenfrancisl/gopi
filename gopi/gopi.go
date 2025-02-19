package gopi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
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
		Paths:      openapi3.NewPaths(),
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

// Marshal the spec into a JSON string
func (gopi *Gopi) MarshalJSON() ([]byte, error) {
	return gopi.spec.MarshalJSON()
}

func (gopi *Gopi) MarshalJSONIndent(prefix string, indent string) ([]byte, error) {
	js, err := gopi.spec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	indented := &bytes.Buffer{}

	err = json.Indent(indented, js, prefix, indent)
	if err != nil {
		return nil, err
	}

	return indented.Bytes(), nil
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

// Initiate the docs for a GET operation on a route
func (route *Route) Get() *Operation {
	pathItem := route.gopi.spec.Paths.Find(route.Path)
	if pathItem == nil {
		return &Operation{}
	}

	pathItem.Get = openapi3.NewOperation()

	return &Operation{
		pathItem: pathItem,
		method:   http.MethodGet,
		route:    route,
	}
}

type Operation struct {
	method   string
	pathItem *openapi3.PathItem
	route    *Route
}

func (op *Operation) getMatchingOperation() *openapi3.Operation {
	if op.method == http.MethodGet {
		return op.pathItem.Get
	}

	if op.method == http.MethodConnect {
		return op.pathItem.Connect
	}

	if op.method == http.MethodDelete {
		return op.pathItem.Delete
	}

	if op.method == http.MethodHead {
		return op.pathItem.Head
	}

	if op.method == http.MethodOptions {
		return op.pathItem.Options
	}

	if op.method == http.MethodPatch {
		return op.pathItem.Patch
	}

	if op.method == http.MethodPost {
		return op.pathItem.Post
	}

	if op.method == http.MethodPut {
		return op.pathItem.Put
	}

	if op.method == http.MethodTrace {
		return op.pathItem.Trace
	}

	return nil
}

// Add tags for the operation
func (op *Operation) Tags(tags []string) *Operation {
	operation := op.getMatchingOperation()
	operation.Tags = tags

	return op
}

// Set the summary of the operation
func (op *Operation) Summary(text string) *Operation {
	operation := op.getMatchingOperation()
	operation.Summary = text

	return op
}

// Set the body of the request
func (op *Operation) Body(model any) *Operation {
	if op.method == http.MethodGet {
		return op
	}

	schemaRef, err := openapi3gen.NewSchemaRefForValue(
		model,
		op.route.gopi.spec.Components.Schemas,
		openapi3gen.CreateComponentSchemas(openapi3gen.ExportComponentSchemasOptions{
			ExportComponentSchemas: true,
			ExportTopLevelSchema:   true,
			ExportGenerics:         true,
		}),
	)
	if err != nil {
		fmt.Println(err)

		return op
	}

	operation := op.getMatchingOperation()
	operation.RequestBody = &openapi3.RequestBodyRef{
		Ref: schemaRef.RefString(),
	}

	return op
}

// Create a new instance of gopi
func New() *Gopi {
	spec := newSpec("test")

	api := &Gopi{
		spec: spec,
	}

	return api
}

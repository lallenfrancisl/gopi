package gopi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/lallenfrancisl/kin-openapi/openapi3"
	"github.com/lallenfrancisl/kin-openapi/openapi3gen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func newSpec(name string) *openapi3.T {
	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:      name,
			Version:    "0.0.0",
			Extensions: map[string]interface{}{},
			Contact:    &openapi3.Contact{},
			License:    &openapi3.License{},
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
	spec      *openapi3.T
	generator *openapi3gen.Generator
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

// Set the title of the api
func (gopi *Gopi) Title(text string) *Gopi {
	gopi.spec.Info.Title = text

	return gopi
}

// Set the description of the api
func (gopi *Gopi) Description(text string) *Gopi {
	gopi.spec.Info.Description = text

	return gopi
}

// Set the terms of service of the api
//
// This must be a valid URL
func (gopi *Gopi) TermsOfService(text string) *Gopi {
	gopi.spec.Info.TermsOfService = text

	return gopi
}

// Struct passed in for defining the contact details
type ContactDef struct {
	Name  string
	URL   string
	Email string
}

// Set the contact details of the api
func (gopi *Gopi) Contact(contact ContactDef) *Gopi {
	gopi.spec.Info.Contact.Name = contact.Name
	gopi.spec.Info.Contact.URL = contact.URL
	gopi.spec.Info.Contact.Email = contact.Email

	return gopi
}

// Struct passed in to set license details of the API
type LicenseDef struct {
	Name string
	URL  string
}

// Set the license details of the API
func (gopi *Gopi) License(license LicenseDef) *Gopi {
	gopi.spec.Info.License.Name = license.Name
	gopi.spec.Info.License.URL = license.URL

	return gopi
}

// Set the version of the api
func (gopi *Gopi) Version(version string) *Gopi {
	gopi.spec.Info.Version = version

	return gopi
}

type ExternalDocDef struct {
	Description string
	URL         string
}

type TagDef struct {
	Name         string
	Description  string
	ExternalDocs ExternalDocDef
}

// Define a tag for use in operation objects
//
// The metadata from this will be shown when the tag is used
// in documentation
func (gopi *Gopi) DefineTag(tag TagDef) *Gopi {
	gopi.spec.Tags = append(gopi.spec.Tags, &openapi3.Tag{
		Name:        tag.Name,
		Description: tag.Description,
		ExternalDocs: &openapi3.ExternalDocs{
			Description: tag.ExternalDocs.Description,
			URL:         tag.ExternalDocs.URL,
		},
	})

	return gopi
}

// Marshal the spec into a JSON string
func (gopi *Gopi) MarshalJSON() ([]byte, error) {
	return gopi.spec.MarshalJSON()
}

// Marshal the spec into a JSON string with indentation
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

func (gopi *Gopi) generateOpenAPISchemaRef(model any, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
	schemaRef, err := gopi.generator.NewSchemaRefForValue(
		model,
		schemas,
	)

	return schemaRef, err
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

var caser = cases.Title(language.English)

func convertToOperationID(input string) string {
	words := strings.Split(input, " ")
	key := strings.ToLower(words[0])

	for _, word := range words[1:] {
		key += caser.String(word)
	}

	return key
}

// Initiate the docs for a GET operation on a route
func (route *Route) Get() *Operation {
	pathItem := route.gopi.spec.Paths.Find(route.Path)
	if pathItem == nil {
		return &Operation{}
	}

	pathItem.Get = openapi3.NewOperation()
	pathItem.Get.Responses = openapi3.NewResponses()

	return &Operation{
		pathItem: pathItem,
		method:   http.MethodGet,
		route:    route,
	}
}

// Add docs for a POST operation on a route
func (route *Route) Post() *Operation {
	pathItem := route.gopi.spec.Paths.Find(route.Path)
	if pathItem == nil {
		return &Operation{}
	}

	pathItem.Post = openapi3.NewOperation()
	pathItem.Post.Responses = openapi3.NewResponses()

	return &Operation{
		pathItem: pathItem,
		method:   http.MethodPost,
		route:    route,
	}
}

// Add docs for a PUT operation on a route
func (route *Route) Put() *Operation {
	pathItem := route.gopi.spec.Paths.Find(route.Path)
	if pathItem == nil {
		return &Operation{}
	}

	pathItem.Put = openapi3.NewOperation()
	pathItem.Put.Responses = openapi3.NewResponses()

	return &Operation{
		pathItem: pathItem,
		method:   http.MethodPut,
		route:    route,
	}
}

// Add docs for a DELETE operation on a route
func (route *Route) Delete() *Operation {
	pathItem := route.gopi.spec.Paths.Find(route.Path)
	if pathItem == nil {
		return &Operation{}
	}

	pathItem.Delete = openapi3.NewOperation()
	pathItem.Delete.Responses = openapi3.NewResponses()

	return &Operation{
		pathItem: pathItem,
		method:   http.MethodDelete,
		route:    route,
	}
}

// Add docs for a PATCH operation on a route
func (route *Route) Patch() *Operation {
	pathItem := route.gopi.spec.Paths.Find(route.Path)
	if pathItem == nil {
		return &Operation{}
	}

	pathItem.Patch = openapi3.NewOperation()
	pathItem.Patch.Responses = openapi3.NewResponses()

	return &Operation{
		pathItem: pathItem,
		method:   http.MethodPatch,
		route:    route,
	}
}

// Struct for holding methods for managing route operation
type Operation struct {
	method   string
	pathItem *openapi3.PathItem
	route    *Route
}

// Add tags for the operation
func (op *Operation) Tags(tags []string) *Operation {
	operation := op.getMatchingOperation()
	operation.Tags = tags

	return op
}

// Deprecate this operation
//
// This must be used at the last since nothing is allowed
// to chain after this
func (op *Operation) Deprecated() {
	operation := op.getMatchingOperation()
	operation.Deprecated = true
}

// Set the summary of the operation
func (op *Operation) Summary(text string) *Operation {
	operation := op.getMatchingOperation()
	operation.Summary = text
	operation.OperationID = convertToOperationID(text)

	return op
}

// Set the description of the operation
func (op *Operation) Description(text string) *Operation {
	operation := op.getMatchingOperation()
	operation.Description = text

	return op
}

// Add documentation for request body
func (op *Operation) Body(model any) *Operation {
	if op.method == http.MethodGet {
		return op
	}

	schemaRef, err := op.route.gopi.generateOpenAPISchemaRef(
		model,
		op.route.gopi.spec.Components.Schemas,
	)

	if err != nil {
		fmt.Println(err)

		return op
	}

	content := openapi3.NewContent()

	contentType := getContentType(model)

	content[contentType] = &openapi3.MediaType{
		Schema: schemaRef,
	}

	operation := op.getMatchingOperation()
	operation.RequestBody = &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Content: content,
		},
	}

	return op
}

// Add documentation for response body
func (op *Operation) Response(status int, model any) *Operation {
	operation := op.getMatchingOperation()
	res := openapi3.NewResponse()

	schemaRef, err := op.route.gopi.generateOpenAPISchemaRef(
		model,
		op.route.gopi.spec.Components.Schemas,
	)
	if err != nil {
		fmt.Println(err.Error())

		return op
	}

	res.WithJSONSchemaRef(schemaRef)
	operation.AddResponse(status, res)

	return op
}

// Struct for holding a route parameter methods
// such as path param, query param etc
type RouteParam struct {
	spec  *openapi3.Parameter
	model any
}

// Create path parameter RouteParam instance
func PathParam(name string, model any) *RouteParam {
	param := &RouteParam{
		model: model,
		spec: &openapi3.Parameter{
			Name:   name,
			Schema: openapi3.NewStringSchema().NewRef(),
			In:     openapi3.ParameterInPath,
		},
	}

	return param
}

// Create query parameter RouteParam instance
func QueryParam(name string, model any) *RouteParam {
	param := &RouteParam{
		model: model,
		spec: &openapi3.Parameter{
			Name:   name,
			Schema: openapi3.NewStringSchema().NewRef(),
			In:     openapi3.ParameterInQuery,
		},
	}

	return param
}

// Create cookie parameter RouteParam instance
func CookieParam(name string, model any) *RouteParam {
	param := &RouteParam{
		model: model,
		spec: &openapi3.Parameter{
			Name:   name,
			Schema: openapi3.NewStringSchema().NewRef(),
			In:     openapi3.ParameterInCookie,
		},
	}

	return param
}

// Create header parameter RouteParam instance
func HeaderParam(name string, model any) *RouteParam {
	param := &RouteParam{
		model: model,
		spec: &openapi3.Parameter{
			Name:   name,
			Schema: openapi3.NewStringSchema().NewRef(),
			In:     openapi3.ParameterInHeader,
		},
	}

	return param
}

// Set the description of the route parameter
func (param *RouteParam) Description(input string) *RouteParam {
	param.spec.Description = input

	return param
}

// Set the param as required
func (param *RouteParam) Required() *RouteParam {
	param.spec.Required = true

	return param
}

// Deprecate the route parameter
//
// This must be specified last, as it doesn't allow further chaining
func (param *RouteParam) Deprecated() {
	param.spec.Deprecated = true
}

// Allow empty values for this parameter
func (param *RouteParam) AllowEmpty() *RouteParam {
	param.spec.AllowEmptyValue = true

	return param
}

// Add RouteParams to an operation
func (op *Operation) Params(params ...*RouteParam) *Operation {
	operation := op.getMatchingOperation()

	for _, param := range params {
		schema, err := op.route.gopi.generateOpenAPISchemaRef(param.model, nil)
		if err != nil {
			return op
		}

		param.spec.Schema = schema
		operation.AddParameter(param.spec)
	}

	return op
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

func getKind(input any) reflect.Kind {
	rv := reflect.ValueOf(input)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	kind := rv.Type().Kind()

	return kind
}

// Get the content type of the go struct
func getContentType(model any) string {
	kind := getKind(model)

	jsonKinds := []reflect.Kind{
		reflect.Struct,
		reflect.Slice,
		reflect.Map,
		reflect.Array,
		reflect.Map,
	}

	if slices.Contains(jsonKinds, kind) {
		return "application/json"
	}

	// Return default fallback content type
	return "application/octet-stream"
}

func splitRefPath(path string) []string {
	if !strings.HasPrefix(path, "#/") {
		return []string{}
	}

	return strings.Split(path, "/")[1:]
}

func newGenerator() *openapi3gen.Generator {
	return openapi3gen.NewGenerator(
		openapi3gen.CreateComponentSchemas(openapi3gen.ExportComponentSchemasOptions{
			ExportComponentSchemas: true,
			ExportTopLevelSchema:   true,
			ExportGenerics:         true,
		}),
		openapi3gen.UseAllExportedFields(),
	)
}

// Create a new instance of gopi
func New() *Gopi {
	spec := newSpec("test")

	api := &Gopi{
		spec:      spec,
		generator: newGenerator(),
	}

	return api
}

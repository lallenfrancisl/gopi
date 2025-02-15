package main

import (
	"fmt"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/index"
	"gopkg.in/yaml.v3"
)

type Node = yaml.Node
type DOM = v3.Document
type Info = base.Info
type Operation struct {
	v3.Operation
}
type Server = v3.Server
type Paths = v3.Paths
type Components = v3.Components
type SecurityRequirement = base.SecurityRequirement
type Tag = base.Tag
type ExternalDoc = base.ExternalDoc
type PathItem = v3.PathItem
type SpecIndex = index.SpecIndex
type Rolodex = index.Rolodex
type Contact = base.Contact
type License = base.License
type ServerVariable = v3.ServerVariable
type Parameter = v3.Parameter
type Response = v3.Response
type Header = v3.Header
type Example = base.Example
type MediaType = v3.MediaType

type API struct {
	dom *DOM
}

type CreateRouteParams struct {
	path        string
	description string
	summary     string
}

func (api *API) Route(params CreateRouteParams) *Route {
	return &Route{
		api:  api,
		path: params.path,
		pathItem: &PathItem{
			Description: params.description,
			Summary:     params.summary,
		},
	}
}

type Route struct {
	api      *API
	path     string
}

func (r *Route) Get(operation Operation) *Route {
	if len(operation.OperationId) == 0 {
		operation.OperationId = fmt.Sprintf("GET %s", r.path)
	}

	r.api.dom.Paths.PathItems.Set(r.path)

	r.pathItem.Get = &operation.Operation

	return r
}

func NewApi(title, version string) *API {
	newAPI := API{
		dom: &DOM{
			Version: "3.0.0",
			Info: &Info{
				Title:   title,
				Version: version,
			},
		},
	}

	return &newAPI
}

package main

import (
	"fmt"
)

func main() {
	api := NewApi("greenlight", "1.0.0")

	fmt.Printf("openapi: %s\n", api.DOM.Version)
	fmt.Printf("%s %s\n", api.DOM.Info.Title, api.DOM.Info.Version)
}

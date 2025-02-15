package main

import (
	"fmt"
)

func main() {
	api := NewApi("greenlight", "1.0.0")

	fmt.Printf("openapi: %s\n", api.dom.Version)
	fmt.Printf("%s %s\n", api.dom.Info.Title, api.dom.Info.Version)
}

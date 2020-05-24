package main

import (
	"os"

	"github.com/chtisgit/go-flows/flows"
	_ "github.com/chtisgit/go-flows/packet" //initialize features
)

func init() {
	addCommand("features", "List available features", listFeatures)
}

func listFeatures(string, []string) {
	//TODO add some kind of limit and filters
	flows.ListFeatures(os.Stdout)
}

package main

import (
	"os"

	"github.com/chtisgit/go-flows/packet"
)

func init() {
	addCommand("keys", "List available keys", listKeys)
}

func listKeys(string, []string) {
	//TODO add some kind of limit and filters
	packet.ListKeys(os.Stdout)
}

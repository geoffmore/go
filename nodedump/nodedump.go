package main

import (
	"flag"
	"node"
)

func main() {
	var node1string, node2string string

	// https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go
	// Boolean flags are true if set and false otherwise
	node.InitFlags()

	// https://stackoverflow.com/questions/19617229/golang-flag-gets-interpreted-as-first-os-args-argument
	// Parse non-command line arguments
	// flag.Args() does not index past the last argument, but flag.Arg(n) returns an empty string
	node1string = flag.Arg(0)
	node2string = flag.Arg(1)
	node1, node2 := node.InitNodes(node1string, node2string)

	// Generate the struct we need to compare nodes
	if node2string != "" {
		node1.Diff(node2)
	} else {
		node1.Dump()
	}
}

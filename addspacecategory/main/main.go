package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/addspacecategory"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a string")
	flag.Parse()
	addspacecategory.AddSpaceCategory(*propPtr)

}

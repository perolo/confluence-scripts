package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/removespacecategory"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a string")
	flag.Parse()
	removespacecategory.RemoveSpaceCategory(*propPtr)

}

package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/addspacecategory"
)

func main() {
	propPtr := flag.String("prop", "jira.properties", "a string")
	flag.Parse()
	addspacecategory.AddSpaceCategory(*propPtr)

}

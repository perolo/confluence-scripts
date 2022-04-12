package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/grouppermissionsreport"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a string")

	grouppermissionsreport.SpacePermissionsReport(*propPtr)
}

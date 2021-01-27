package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/spacepermissionsreport"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a string")

	spacepermissionsreport.SpacePermissionsReport(*propPtr)
}

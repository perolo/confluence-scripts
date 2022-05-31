package main

import (
	"flag"

	"github.com/perolo/confluence-scripts/spacepermissionsreport"
)

func main() {
	propPtr := flag.String("prop", "spacepermissionsreport.properties", "a string")

	spacepermissionsreport.SpacePermissionsReport(*propPtr)
}

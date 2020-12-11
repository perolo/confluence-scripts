package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/syncadgroup"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a properties file")

	syncadgroup.ConfluenceSyncAdGroup(*propPtr)
}

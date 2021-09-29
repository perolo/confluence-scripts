package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/synccofluenceadgroup"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a properties file")

	synccofluenceadgroup.ConfluenceSyncAdGroup(*propPtr)
}

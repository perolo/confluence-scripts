package main

import (
	"flag"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-scripts/spacepermissionsmodifier"
	"log"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	// or through Decode
	var cfg spacepermissionsmodifier.SpacePermissionsModConfig
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	spacepermissionsmodifier.Modify(cfg)
}

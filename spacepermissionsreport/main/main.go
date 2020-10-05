package main

import (
	"flag"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-scripts/spacepermissionsreport"
	"log"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	// or through Decode
	var cfg spacepermissionsreport.SpacePermissionsReportConfig
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

//	cfg.File = "C:\\Users\\perolo\\Documents\\Tmp.xlsx"
//	cfg.Space = ""
//	cfg.SpaceCategory = cfg.SpaceCategory

	spacepermissionsreport.SpacePermissionsReport(cfg)
}

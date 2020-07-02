package main

import (
	"github.com/perolo/confluence-prop/client"
	"flag"
	"github.com/magiconair/properties"
	"log"
	"fmt"
)

func main() {

	propPtr := flag.String("prop", "../confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	// or through Decode
	type Config struct {
		ConfHost    string `properties:"confhost"`
		User        string `properties:"user"`
		Pass        string `properties:"password"`
		Source      string `properties:"source"`
		Destination string `properties:"destination"`
	}
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	theClient := client.Client(&config)

	start := 0
	cont := true
	increase := 10

	noSpaces := 0
	for cont {
		opt := client.GroupOptions{Start: start, Limit: increase}
		spaces := theClient.GetAllSpacesForGroupPermissions(cfg.Source, &opt)

		for _, space := range spaces.Spaces {
			fmt.Printf("Space name: %s\n", space.Name)
			p := space.Permissions
			added := theClient.AddSpacePermissionsForGroup(space.Key, cfg.Destination, p)
			fmt.Printf("Permissions added : %s\n", added.Added)
			noSpaces++
		}

		start = start + increase
		if spaces.Total < increase {
			cont = false
		}
	}
	fmt.Printf("Spaces : %d \n", noSpaces)

}
